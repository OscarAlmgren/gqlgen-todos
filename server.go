package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/oscaralmgren/hackernews/graph"
	"github.com/oscaralmgren/hackernews/graph/generated"
	"github.com/oscaralmgren/hackernews/internal/auth"
	mongodb "github.com/oscaralmgren/hackernews/internal/pkg/db/migrations/mongodb"
	mysqldb "github.com/oscaralmgren/hackernews/internal/pkg/db/migrations/mysql"
)

const defaultPort = "8080"

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(auth.Middleware())

	// init MySQL db, then migrate and close
	mysqldb.InitDB()
	defer mysqldb.CloseDB()
	mysqldb.Migrate()

	mongodb.InitMongoDB()
	mongodb.PingMongoDB()
	defer mongodb.CloseMongoDB()

	server := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", server)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router)) // router handler instead of nil
}
