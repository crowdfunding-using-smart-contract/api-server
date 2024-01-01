.PHONY: run
run:
	go run cmd/api/main.go

.PHONY: swag
swag:
	swag init --parseDependency -g cmd/api/server/server.go