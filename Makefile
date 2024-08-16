.PHONY: build test

build:
	go build ./...

test: bin/sqlc-graphql.wasm
	go test ./...

all: bin/sqlc-graphql bin/sqlc-graphql.wasm

bin/sqlc-graphql: bin go.mod go.sum $(wildcard **/*.go)
	cd plugin && go build -o ../bin/sqlc-graphql ./main.go

bin/sqlc-graphql.wasm: bin/sqlc-graphql
	cd plugin && GOOS=wasip1 GOARCH=wasm go build -o ../bin/sqlc-graphql.wasm main.go

bin:
	mkdir -p bin
