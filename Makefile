.PHONY: run
run:
	go run cmd/api/main.go

.PHONY: swag
swag:
	swag init --parseDependency -g cmd/api/server/server.go

.PHONY: test
test:
	go test -v ./...