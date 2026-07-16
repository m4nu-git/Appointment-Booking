.PHONY: run dev build test test-integration lint \
        docker-up docker-down docker-build \
        migrate-up migrate-down migrate-create

# --- Local (no Docker) ------------------------------------------------------

run:
	go run ./cmd/api

dev:
	air -c .air.toml

build:
	go build -o bin/api ./cmd/api

test:
	go test ./... -race -cover

test-integration:
	go test ./test/integration/... -race -tags=integration

lint:
	golangci-lint run

# --- Docker ------------------------------------------------------------------

docker-up:
	docker compose up --build

docker-down:
	docker compose down

docker-build:
	docker build -t appointment-booking-api .

# --- Migrations (Phase 1+) ----------------------------------------------------

DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@localhost:$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

migrate-up:
	migrate -path migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path migrations -database "$(DB_URL)" down

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)