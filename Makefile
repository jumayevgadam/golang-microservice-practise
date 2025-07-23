.PHONY: help build-all build run test clean lint generate-stocks-proto generate-stocks-for-cart generate-cart-proto

# Default target
help:
	@echo "Available targets:"
	@echo "  build-all   - Build all services for Linux (amd64)"
	@echo "  build       - Build all services for current OS"
	@echo "  run         - Run all services locally"
	@echo "  test        - Run tests"
	@echo "  lint        - Run golangci-lint across all services"
	@echo "  clean       - Remove build artifacts"

# Cross-platform build for Linux (e.g., for deployment)
build-all:
	@echo "Building all services for Linux (amd64)..."
	@cd cart && GOOS=linux GOARCH=amd64 $(MAKE) build
	@cd stocks && GOOS=linux GOARCH=amd64 $(MAKE) build

# Local development build (current OS)
build:
	@echo "Building all services for $(shell uname -s)/$(shell uname -m)..."
	@$(MAKE) -C cart build
	@$(MAKE) -C stocks build

run:
	@echo "Running services (logs will show below)..."
	@echo "=== Cart Service ==="
	@$(MAKE) -C cart run
	@echo "=== Stocks Service ==="
	@$(MAKE) -C stocks run

lint:
	@echo "Running golangci-lint…"
	# Point at each module directory, or simply `./…` if you want everything
	golangci-lint run ./cart/... ./stocks/...


test:
	@$(MAKE) -C cart test
	@$(MAKE) -C stocks test

clean:
	@$(MAKE) -C cart clean
	@$(MAKE) -C stocks clean

generate-stocks-proto:
	@mkdir -p stocks/pkg/api/stocks
	@protoc \
		-I proto \
		-I ./stocks/vendor.protogen \
		--go_out=stocks/pkg/api/stocks --go_opt=paths=source_relative \
		--go-grpc_out=stocks/pkg/api/stocks --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=stocks/pkg/api/stocks --grpc-gateway_opt=paths=source_relative \
		proto/stocks.proto

generate-stocks-for-cart:
	@mkdir -p cart/pkg/api/stocks
	@protoc \
		-I proto \
		-I ./stocks/vendor.protogen \
		--go_out=cart/pkg/api/stocks --go_opt=paths=source_relative \
		--go-grpc_out=cart/pkg/api/stocks --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=cart/pkg/api/stocks --grpc-gateway_opt=paths=source_relative \
		proto/stocks.proto

generate-cart-proto:
	@mkdir -p cart/pkg/api/cart
	@protoc \
		-I proto \
		-I ./cart/vendor.protogen \
		--go_out=cart/pkg/api/cart --go_opt=paths=source_relative \
		--go-grpc_out=cart/pkg/api/cart --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=cart/pkg/api/cart --grpc-gateway_opt=paths=source_relative \
		proto/cart.proto

		