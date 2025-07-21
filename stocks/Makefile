APP_NAME := stocks
BIN_DIR := bin
PORT := 8081
MIGRATION_DIR := $(CURDIR)/internal/migrations

ifeq ($(POSTGRES_SETUP),)
	POSTGRES_SETUP := user=postgres password=12345 dbname=stock_service_db host=localhost port=5432 sslmode=disable
endif
LOCAL_BIN=$(CURDIR)/bin

.PHONY: build run test clean migration-create migration-up migration-down migration-status migration-reset generate-mock docker-build install-deps vendor-proto get-deps generate-abstract-proto generate-stocks-proto

build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/$(APP_NAME) ./cmd/main.go

docker-build:
	@docker compose up -d --build

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

migration-create:
	@mkdir -p $(MIGRATION_DIR)
	@goose -dir $(MIGRATION_DIR) create $(name) sql

migration-up:
	@goose -dir $(MIGRATION_DIR) postgres "$(POSTGRES_SETUP)" up

migration-down:
	@goose -dir $(MIGRATION_DIR) postgres "$(POSTGRES_SETUP)" down

migration-reset:
	@goose -dir $(MIGRATION_DIR) postgres "$(POSTGRES_SETUP)" reset

generate-mock:
	@go generate ./...

install-deps:
	@mkdir -p $(LOCAL_BIN)
	@GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	@GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

vendor-proto:
	@if [ ! -d vendor.protogen/google ]; then \
		git clone https://github.com/googleapis/googleapis vendor.protogen/googleapis && \
		mkdir -p vendor.protogen/google/ && \
		mv vendor.protogen/googleapis/google/api vendor.protogen/google && \
		rm -rf vendor.protogen/googleapis; \
	fi

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go get -u github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway

generate-abstract-proto: 
	@mkdir -p pkg/api/abstract
	@protoc --proto_path=api/abstract --proto_path=vendor.protogen \
 		--go_out=pkg/api/abstract --go_opt=paths=source_relative \
 		--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go \
 		api/abstract/abstract.proto

generate-stocks-proto:
	@mkdir -p pkg/api/stocks
	@protoc --proto_path=api/stocks --proto_path=vendor.protogen \
		--go_out=pkg/api --go_opt=paths=source_relative \
		--go-grpc_out=pkg/api --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=pkg/api --grpc-gateway_opt=paths=source_relative \
		api/stocks/stocks_v1.proto