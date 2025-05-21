.PHONY: build run test migrate-up migrate-down swagger lint clean

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=rota-api

# Load environment variables
include .env

export DB_DSN=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

run:
	$(GOCMD) run main.go

swagger:
	swag init -g main.go --output ./docs

migrate-up:
	migrate -path ./migrations -database "$(DB_DSN)" up

migrate-down:
	migrate -path ./migrations -database "$(DB_DSN)" down

lint:
	golangci-lint run

setup-dev:
	# Install development tools
	go install github.com/cosmtrek/air@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/golang/mock/mockgen@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Generate mocks for testing
generate-mocks:
	mockgen -source=./repositories/user_repository.go -destination=./mocks/user_repository_mock.go -package=mocks
	# Add more mocks as needed

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application"
	@echo "  test          - Run tests"
	@echo "  migrate-up    - Run database migrations"
	@echo "  migrate-down  - Rollback database migrations"
	@echo "  swagger       - Generate Swagger documentation"
	@echo "  lint          - Run linter"
	@echo "  setup-dev     - Install development tools"
	@echo "  clean         - Clean build artifacts"
	@echo "  help          - Show this help message"
