.PHONY: build run test clean docker-up docker-down migrate-up migrate-down seed

# Build the application
build:
	go build -o server cmd/server/main.go

# Run the application
run:
	go run cmd/server/main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -f server
	go clean

# Docker commands
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

# Database migrations
migrate-up:
	go run cmd/server/main.go migrate

migrate-down:
	@echo "Down migrations should be handled with care. Use psql manual if needed or implement in Go."

# Seed data
seed:
	psql -h localhost -U postgres -d tenangantri -f migrations/002_seed_data.up.sql

# Development setup
dev-setup:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Lint
lint:
	golangci-lint run

# Generate mocks (if needed)
generate:
	go generate ./...
