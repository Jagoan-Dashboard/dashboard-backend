# Build stage
FROM golang:1.24-alpine AS builder

# Install dependencies
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

WORKDIR /app

# Copy modules and download dependencies
COPY go.mod go.sum ./ 
RUN go mod download
RUN go mod verify

# Copy all source code
COPY . .

# Install goose binary
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Build API binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o main cmd/api/main.go

# Build Seeder binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o seeder cmd/seeder/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

# Create a non-root user
RUN adduser -D -g '' appuser

# Copy migrations files
COPY --from=builder /app/migrations /app/migrations

# Copy binaries from builder stage
COPY --from=builder /go/bin/goose /usr/local/bin/
COPY --from=builder /app/main /app/main
COPY --from=builder /app/seeder /app/seeder

USER appuser

EXPOSE 8080

# Default command to run the API
ENTRYPOINT ["/app/main"]
