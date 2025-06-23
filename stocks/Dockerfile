FROM golang:1.24-alpine3.21 AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

# Copy Source code
COPY .  .

# Build Cart service app
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o bin/stocks ./cmd/main.go

# Final stage
FROM alpine:3.21

WORKDIR /root

# Copy binary from builder stage
COPY --from=builder /app/bin/stocks .

COPY .env .

EXPOSE 8081

CMD [ "./stocks" ]