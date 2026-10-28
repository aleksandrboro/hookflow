.DEFAULT_GOAL := start
.PHONY: start test up-migrations down-migrations

start:
	go run cmd/hookflow/main.go --config=config/local.yaml

test:
	go test -v -race -timeout 30s ./...

up-migrations:
	go run cmd/migrator/migrator.go --migrations-path=./migrations --migrations-type=up

down-migrations:
	go run cmd/migrator/migrator.go --migrations-path=./migrations --migrations-type=down