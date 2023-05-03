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
	"strconv"
	"strings"
	"time"
)

const defaultHttpPort = 8080
const defaultGrpcPort = 5050

type Config struct {
	// gorm postgres connection url
	postgresConnectionUrl string
	// redis connection string "redis://<user>:<pass>@localhost:6379/<db>"
	redisConnectionString string
	// application http port
	httpPort int
	// application grpc port
	grpcPort int
	// allowed origins to connect to websocket
	allowedOrigins []string
	// allowed core origin
	allowedCoreOrigin string
}

func Configure() (Config, error) {
	var config = Config{}
	config.postgresConnectionUrl = os.Getenv("PG_DSN")
	config.redisConnectionString = os.Getenv("REDIS_CONN_STRING")
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		config.httpPort = defaultHttpPort
	} else {
		port, err := strconv.Atoi(httpPort)
		if err != nil {
			return config, err
		}
		config.httpPort = port
	}
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		config.grpcPort = defaultGrpcPort
	} else {
		port, err := strconv.Atoi(grpcPort)
		if err != nil {
			return config, err
		}
		config.grpcPort = port
	}
	config.allowedOrigins = strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")
	config.allowedCoreOrigin = os.Getenv("ALLOWED_CORE_ORIGIN")
	return config, nil
}

func main() {
	// Collect prerequisites
	config, err := Configure()
	if err != nil {
		panic(err)
	}
	// Initiate Gorm connection
	db, err := gorm.Open(postgres.Open(config.postgresConnectionUrl), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	err = models.MigrateSchema(db)
	if err != nil {
		panic("failed to migrate database")
	}

	// Initiate Redis connection
	opt, err := redis.ParseURL(config.redisConnectionString)
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
	websocketUpgrader := websocket.Upgrader{
		CheckOrigin: utils.CheckAllowedOrigin(config.allowedOrigins, true)}
	srv.AddTransport(&transport.Websocket{
		Upgrader:              websocketUpgrader,
		KeepAlivePingInterval: 10 * time.Second,
	})

	// initiate data loaders
	loaders := dataloader.NewLoaders(db)
	dataloaderSrv := dataloader.Middleware(loaders, srv)

	http.Handle("/graphiql/", playground.Handler("GraphQL playground", "/graphql/"))
	http.Handle("/graphql/", utils.CorsMiddleware(dataloaderSrv, config.allowedCoreOrigin))

	log.Printf("connect to http://localhost:%d/graphiql/ for GraphQL playground", config.httpPort)
	go func() { log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.httpPort), nil)) }()

	// grpc server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterCommentingServer(grpcServer, &grpc_gate.CommentingServer{DB: db, RDB: rdb})

	log.Printf("grpc server listening at %v", lis.Addr())
	go func() { log.Fatal(grpcServer.Serve(lis)) }()

	select {}
}
