.PHONY: start start-db help

help:
	@echo "Available commands:"
	@echo "---------------------"
	@echo "start: Start the application server"
	@echo "start-db: Start the containerized database instance"
	@echo "gen-docs: Generate the OpenAPI swagger documentation"

start:
	go run cmd/main.go

start-db:
	docker compose up -d

gen-docs:
	swag init -g ./cmd/main.go -o ./docs
