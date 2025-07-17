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

## Homework 10

### ✅ Task Overview

- Replacing all HTTP handlers with **gRPC** service definitions.
- Compatible with tools like [grpcui](https://github.com/fullstorydev/grpcui).
- (Bonus) Adding **gRPC-Gateway** support to allow access via both **HTTP/REST** and **gRPC** clients.

### 🧱 Project Structure

```
.
├── pkg/
│   └── api/                     # Generated gRPC & Gateway code
│       ├── service.proto
|       |   service.pb.go
│       ├── service_grpc.pb.go
│       └── service.pb.gw.go
├── internal/
│   ├── service/                 # Business logic
│   └── server/
│       ├── grpc.go              # gRPC server setup
│       └── gateway.go           # gRPC-Gateway HTTP server setup
├── cmd/
│   └── main.go                  # Entrypoint
├── go.mod
└── README.md                    # You are here
```

### 🧪 How to Generate Code from `.proto`

Install protoc plugins if not already:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
```

Then generate the code:

```bash
protoc -I proto \
  --go_out=pkg/api --go_opt=paths=source_relative \
  --go-grpc_out=pkg/api --go-grpc_opt=paths=source_relative \
  --grpc-gateway_out=pkg/api --grpc-gateway_opt=paths=source_relative \
  proto/service.proto
```


### 🧰 Testing the Service

### gRPC Call (CLI)

```bash
grpcurl -plaintext localhost:9090 list
```

### Bonus HTTP Call (REST via gRPC-Gateway)

```bash
curl "http://localhost:8080/v1/data?id=123"
```

### UI Test (gRPC UI)

```bash
grpcui -plaintext localhost:9090
```

Then open the browser at: [http://localhost:8080](http://localhost:9090)

### 🧾 Notes

- All proto definitions and generated Go code are stored in `pkg/api/`.
- This project supports both gRPC and REST clients.
- Fully testable with [grpcurl](https://github.com/fullstorydev/grpcurl) and [grpcui](https://github.com/fullstorydev/grpcui).


### 🐳 Docker Instructions

Make sure to rebuild your Docker images after applying the gRPC and gRPC-Gateway changes:

```bash
# Example Docker build for the service
docker build -t service_name:hw10 .
```

> ⚠️ Don’t forget to update your Dockerfile to install.
