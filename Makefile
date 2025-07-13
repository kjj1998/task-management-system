include .env
export $(shell sed 's/=.*//' .env)

# Start app
run-dev:
		ENV=dev go run cmd/api/main.go

build-dev:
		docker compose build --no-cache

# Start app and db in dev environment:
start-dev:
		docker compose down && docker compose up

# Run migrations
migrate-up:
		migrate -database "mysql://$(DB_USER):$(DB_PASS)@tcp($(DB_HOST):$(DB_PORT))/$(DB)" -path migrations up

migrate-up-once:
		migrate -database "mysql://$(DB_USER):$(DB_PASS)@tcp($(DB_HOST):$(DB_PORT))/$(DB)" -path migrations up 1

migrate-down:
		migrate -database "mysql://$(DB_USER):$(DB_PASS)@tcp($(DB_HOST):$(DB_PORT))/$(DB)" -path migrations down 1

# Run all repository test files
repository-tests:
		go test ./internal/repository/...

lint:
		golangci-lint run
