APP_NAME := cart
BIN_DIR := bin
PORT := 8080
MIGRATION_DIR := $(CURDIR)/internal/migrations

ifeq ($(POSTGRES_SETUP),)
	POSTGRES_SETUP := user=postgres password=12345 dbname=cart_service_db host=localhost port=5432 sslmode=disable
endif

.PHONY: build run test clean migration-create migration-up migration-down migration-reset docker-run

build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/$(APP_NAME) ./cmd/main.go

run: build
	@echo "Running $(APP_NAME) on port $(PORT)..."
	@./$(BIN_DIR)/$(APP_NAME) -port $(PORT)

test:
	@echo "Testing $(APP_NAME)..."
	@go test -v ./...

clean:
	@echo "Cleaning $(APP_NAME)..."
	@rm -rf $(BIN_DIR)
	@go clean

docker-run:
	@docker compose up -d --build

migration-create:
	@mkdir -p $(MIGRATION_DIR)
	@goose -dir $(MIGRATION_DIR) create $(name) sql

migration-up:
	@goose -dir $(MIGRATION_DIR) postgres "$(POSTGRES_SETUP)" up

migration-down:
	@goose -dir $(MIGRATION_DIR) postgres "$(POSTGRES_SETUP)" down

migration-reset:
	@goose -dir $(MIGRATION_DIR) postgres "$(POSTGRES_SETUP)" reset
