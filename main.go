package main

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/SayIfOrg/say_keeper/graph"
	"github.com/SayIfOrg/say_keeper/models"
	"github.com/SayIfOrg/say_keeper/utils"
	"github.com/gorilla/websocket"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const defaultPort = "8080"

func main() {
	// Collect prerequisites
	dsn := "host=localhost user=postgres password=password dbname=keeper port=5432 sslmode=disable"
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
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

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{DB: db}}))

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

	http.Handle("/graphiql/", playground.Handler("GraphQL playground", "/graphql/"))
	http.Handle("/graphql/", utils.CorsMiddleware(srv, os.Getenv("ALLOWED_CORE_ORIGIN")))

	log.Printf("connect to http://localhost:%s/graphiql/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
