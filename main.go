package main

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/SayIfOrg/say_keeper/commenting"
	"github.com/SayIfOrg/say_keeper/gateway/grpc_gate"
	"github.com/SayIfOrg/say_keeper/graph"
	"github.com/SayIfOrg/say_keeper/graph/dataloader"
	"github.com/SayIfOrg/say_keeper/graph/gmodel"
	"github.com/SayIfOrg/say_keeper/models"
	"github.com/SayIfOrg/say_keeper/utils"
	pb "github.com/SayIfOrg/say_protos/packages/go"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

const defaultHttpPort = "8080"
const defaultGrpcPort = "5050"

func main() {
	// Collect prerequisites
	// gorm postgres connection string
	dsn := "host=localhost user=postgres password=password dbname=keeper port=5432 sslmode=disable"
	// redis connection string "redis://<user>:<pass>@localhost:6379/<db>"
	redisURL := "redis://@172.19.97.252:6378/5"
	// application http port
	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = defaultHttpPort
	}
	// application grpc port
	grpcPort := os.Getenv("PORT")
	if grpcPort == "" {
		grpcPort = defaultGrpcPort
	}

	// Initiate Gorm connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	err = models.MigrateSchema(db)
	if err != nil {
		panic("failed to migrate database")
	}

	// Initiate Redis connection
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		panic(err)
	}
	rdb := redis.NewClient(opt)

	ctx := context.Background()

	var subs = new(commenting.Subs[gmodel.Comment])
	// Initiate Redis pub/sub to the comments chan
	go commenting.SubscribeComment[gmodel.Comment](ctx, rdb, subs, gmodel.UnmarshalComment)

	srv := handler.New(graph.NewExecutableSchema(
		graph.Config{Resolvers: &graph.Resolver{DB: db, RDB: rdb, Subs: subs}}))

	// Add ServerSentEvent transport (order should be before transport.POST)
	srv.AddTransport(transport.SSE{})

	// Some default http behaviors
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})

	srv.SetQueryCache(lru.New(1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(100),
	})

	// Add websocket transport
	allowedOrigins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")
	websocketUpgrader := websocket.Upgrader{CheckOrigin: utils.CheckAllowedOrigin(allowedOrigins, true)}
	srv.AddTransport(&transport.Websocket{
		Upgrader:              websocketUpgrader,
		KeepAlivePingInterval: 10 * time.Second,
	})

	// initiate data loaders
	loaders := dataloader.NewLoaders(db)
	dataloaderSrv := dataloader.Middleware(loaders, srv)

	http.Handle("/graphiql/", playground.Handler("GraphQL playground", "/graphql/"))
	http.Handle("/graphql/", utils.CorsMiddleware(dataloaderSrv, os.Getenv("ALLOWED_CORE_ORIGIN")))

	log.Printf("connect to http://localhost:%s/graphiql/ for GraphQL playground", httpPort)
	go func() { log.Fatal(http.ListenAndServe(":"+httpPort, nil)) }()

	// grpc server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterCommentingServer(grpcServer, &grpc_gate.CommentingServer{DB: db, RDB: rdb})

	log.Printf("grpc server listening at %v", lis.Addr())
	go func() {
		log.Fatal(grpcServer.Serve(lis))
	}()

	select {}
}
