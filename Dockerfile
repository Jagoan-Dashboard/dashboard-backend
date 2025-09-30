# Build stage
FROM golang:1.24-alpine AS builder

# Install git for private repos if needed
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

# Create appuser
RUN adduser -D -g '' appuser

WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o main cmd/api/main.go

# Final stage - use alpine for better compatibility
FROM alpine:latest

# Install ca-certificates for HTTPS requests and wget for healthcheck
RUN apk --no-cache add ca-certificates wget tzdata

# Create non-root user
RUN adduser -D -g '' appuser

# Copy the binary
COPY --from=builder /app/main /app/main

# Use non-root user
USER appuser

EXPOSE 8080

ENTRYPOINT ["/app/main"]