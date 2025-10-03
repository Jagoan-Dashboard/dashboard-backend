# Build stage
FROM golang:1.24-alpine AS builder

RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify

COPY . .

# Build API binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o main cmd/api/main.go

# Build Seeder binary (TAMBAHKAN INI)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o seeder cmd/seeder/main.go
    
# Final stage - use alpine for better compatibility
FROM alpine:latest

RUN apk --no-cache add ca-certificates wget tzdata

RUN adduser -D -g '' appuser

# Copy KEDUA binary dari builder stage
COPY --from=builder /app/main /app/main
COPY --from=builder /app/seeder /app/seeder  

USER appuser

EXPOSE 8081

ENTRYPOINT ["/app/main"]