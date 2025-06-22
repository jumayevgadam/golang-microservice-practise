.PHONY: help build-all build run test clean lint

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
