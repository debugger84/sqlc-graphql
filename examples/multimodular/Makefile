include .env

help: ## Commands list
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-16s\033[0m %s\n", $$1, $$2}'

.PHONY: migrate
migrate: ## Run migrations in a database
	docker run --rm -it --network=host -v "$(shell pwd)/comment/storage:/db" ghcr.io/amacneil/dbmate -u ${PG_URI} -d "./db/migration" --no-dump-schema up
	docker run --rm -it --network=host -v "$(shell pwd)/post/storage:/db" ghcr.io/amacneil/dbmate -u ${PG_URI} -d "./db/migration" --no-dump-schema up

.PHONY: rollback
rollback: ## Roll back the latest migration in a database
	echo $(pwd)
	docker run --rm -it --network=host -v "$(shell pwd)/comment/storage:/db" ghcr.io/amacneil/dbmate -u ${PG_URI} -d "./db/migration" --no-dump-schema rollback
	docker run --rm -it --network=host -v "$(shell pwd)/post/storage:/db" ghcr.io/amacneil/dbmate -u ${PG_URI} -d "./db/migration" --no-dump-schema rollback


.PHONY: generate
generate: ## Generate DTO and DAO for modules. After that, generate public graphql schema
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	find . -path './*/*/sqlc.yaml' -exec sqlc -f '{}' generate ';'
	go run github.com/99designs/gqlgen generate --config gqlgen.yml

.PHONY: install
install: ## Make a binary to ./bin folder
	go build -o ./bin/server  ./server.go

.PHONY: run
run: ## Run the application
	$(MAKE) install
	./bin/server