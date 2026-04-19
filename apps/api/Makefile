.PHONY: run build test tidy sqlc migrate-up migrate-down migrate-status lint

# Build the server binary
build:
	go build -o ../../bin/api ./cmd/server

# Run the server (requires .env)
run:
	go run ./cmd/server

# Run all tests
test:
	go test ./... -race -count=1

# Tidy modules
tidy:
	go mod tidy

# Regenerate sqlc types + queries
sqlc:
	sqlc generate

# Run migrations (set DATABASE_URL in env or .env)
migrate-up:
	goose -dir migrations postgres "$(DATABASE_URL)" up

migrate-down:
	goose -dir migrations postgres "$(DATABASE_URL)" down

migrate-status:
	goose -dir migrations postgres "$(DATABASE_URL)" status

# Lint (requires golangci-lint)
lint:
	golangci-lint run ./...
