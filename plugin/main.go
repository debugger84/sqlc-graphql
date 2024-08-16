package main

import (
	"github.com/sqlc-dev/plugin-sdk-go/codegen"
	"os"

	golang "github.com/debugger84/sqlc-graphql/internal"
)

func main() {
	restoreFile := os.Getenv("RESTORE_REQUEST_FILE")

	if restoreFile != "" {
		fd, err := os.Open(restoreFile)
		if err != nil {
			panic(err)
		}
		os.Stdin = fd
	}

	codegen.Run(golang.Generate)
}
