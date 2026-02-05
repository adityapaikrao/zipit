# Stage 1: Build
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

WORKDIR /app

# Copy go mod and sum files from the backend directory
# (Assuming build context is set to the backend folder)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o url-service ./cmd/url-service/main.go

# Stage 2: Runtime
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/url-service .

# The url-service usually runs on a private gRPC port
EXPOSE 5051

# Command to run the service
CMD ["./url-service"]
