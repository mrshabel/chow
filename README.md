# CHOW: Local Food Compass

Community-driven platform built to help you discover local food hot-spots closer to you. Catch up with highly rated joints powered by locals like you.

## Setup

1. Set environment and start the database with docker. A postgres instance with postgis installed is started

```bash
# setup env
cp .env.example .env

# start database container
docker compose up db -d
```

The database instance starts on the port specified in the env file

2. Run the server:

```bash
go run cmd/main.go
```

The server starts on `localhost:8000` by default

## Usage

## TODO
