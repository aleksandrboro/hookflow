.PHONY:
start:
	go run cmd/main.go -config=config/local.yaml

.PHONY:
test:
	go test -v -race -timeout 30s ./...

.DEFAULT_GOAL: start