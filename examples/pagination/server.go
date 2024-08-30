package main

import (
	"context"
	"github.com/jackc/pgx/v5"
	"log"
	"net/http"
	"os"
	"pagination/graph/gen"
	"pagination/graph/resolver"
	"pagination/storage"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "8080"
const defaultPgURI = "postgres://postgres:foobar@127.0.0.1:5432/test?sslmode=disable"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	pgUri := os.Getenv("PG_URI")
	if pgUri == "" {
		pgUri = defaultPgURI
	}

	ctx := context.Background()

	conn, err := pgx.Connect(ctx, pgUri)
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer conn.Close(ctx)

	queries := storage.New(conn)

	srv := handler.NewDefaultServer(gen.NewExecutableSchema(gen.Config{Resolvers: &resolver.Resolver{Queries: queries}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
