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

## Homework 7

### 1. Write a Dockerfile

- Must build and run your application (e.g., cart, stock, etc.)

- Should expose the correct internal port via EXPOSE

- Must define CMD or ENTRYPOINT to start the app

### 2. Write a docker-compose.yml

- Should run your app and any required dependencies (e.g., PostgreSQL)

- All containers must start and communicate correctly

### 3. Push your app's Docker image to a public registry

- Use Docker Hub or another accessible container registry

- Tag the image like:
```bash
    docker tag my-app username/my-app:hw7
    docker push username/my-app:hw7
```

I should be able to pull your image using:
```bash
    docker pull username/my-app:hw7
```

### 4. Document the following in every service's README:

- The Docker image name & tag

- App port (e.g. 8080)

- Required environment variables (e.g., DB_HOST, DB_PORT, etc.)

- Sample requests or endpoints if available


## Homework 8

### Requirements:
- Cover the handlers and use cases with unit tests. Minimum coverage: 40%. 
- Cover the handlers with integration tests. Minimum test cases: successful execution and (receiving an error due to invalid input data or receiving not found error).
```bash
  INTEGRATION_TEST=1
```
Prepare a Makefile for each service that includes the following commands:

  - starting the test environment using docker-compose, 
  - running integration tests, 
  - running unit tests.

After completing all changes, don’t forget to update your Docker Hub images.


## Homework 9
- [Kafka-service](metrics-consumer/README.md)