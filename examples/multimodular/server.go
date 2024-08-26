package main

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/subosito/gotenv"
	"log"
	commentStorage "multimodular/comment/storage"
	"multimodular/graph/gen"
	"multimodular/graph/resolver"
	"multimodular/post/storage"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "8080"
const defaultPgURI = "postgres://postgres:foobar@127.0.0.1:5432/test?sslmode=disable"

func main() {
	loadEnv(".env")
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

	postQueries := storage.New(conn)
	commentQueries := commentStorage.New(conn)

	srv := handler.NewDefaultServer(
		gen.NewExecutableSchema(
			gen.Config{
				Resolvers: &resolver.Resolver{
					PostQueries:    postQueries,
					CommentQueries: commentQueries,
				},
			},
		),
	)

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func loadEnv(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return
		}
		panic(err)
	}

	err = gotenv.Apply(f)

	f.Close()
	if err != nil {
		panic(err)
	}
}
