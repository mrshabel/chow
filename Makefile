.PHONY: start start-db help

help:
	@echo "Available commands:"
	@echo "---------------------"
	@echo "start: Start the application server"
	@echo "start-db: Start the containerized database instance"

start:
	go run cmd/main.go

start-db:
	docker compose up -d
