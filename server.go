package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/oscaralmgren/hackernews/graph"
	"github.com/oscaralmgren/hackernews/graph/generated"
	mysqldb "github.com/oscaralmgren/hackernews/internal/pkg/db/migrations/mysql"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// router := chi.NewRouter()

	// init MySQL db, then migrate and close
	mysqldb.InitDB()
	defer mysqldb.CloseDB()
	mysqldb.Migrate()

	server := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", server)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
