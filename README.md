# Cart and stock Services

## Overview

This repository contains two microservices implemented in Go:
- **Cart Service** (`cart`): Manages shopping cart operations.
- **Stocks Service** (`stocks`): Manages how many items do we have and it's price.

Both services follow a layered (clean/hexagonal) architecture and communicate via net/http package.


## Getting Started

### Build

```bash
make build
```

### Run All Services

```bash
make run
```
Starts `cart`, and `stocks` services.

### Run Linter Services

```bash
make lint
```

### Service Endpoints

Both services expose HTTP endpoints on port 8080 (`cart`) and 8081 (`stocks`).


### Inter-Service Communication

- Uses HTTP

## Testing

- Unit tests for core business logic:
  ```bash
  make test 
  ```

## Documentation

Create service skeletons for cart and stocks services according to the [documentation](docs/README.md).

---
