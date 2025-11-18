# Multi-stage Dockerfile for WebSocket Echo Server

# Stage 1: Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /build

# Copy go module files first (for better caching)
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY websocket-echo-server.go .

# Build the application
# CGO_ENABLED=0 for static binary
# -ldflags="-w -s" to reduce binary size
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o echo-server \
    websocket-echo-server.go

# Expose port
EXPOSE 8080

# Run the application
ENTRYPOINT ["./echo-server"]
