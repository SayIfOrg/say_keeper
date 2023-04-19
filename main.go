package main

import (
	"github.com/SayIfOrg/say_keeper/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/SayIfOrg/say_keeper/graph"
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

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{DB: db}}))

	http.Handle("/graphiql/", playground.Handler("GraphQL playground", "/graphql/"))
	http.Handle("/graphql/", srv)

	log.Printf("connect to http://localhost:%s/graphiql/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
