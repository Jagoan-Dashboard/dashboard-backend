# Makefile
.PHONY: help run build test migrate

help:
	@echo "Available commands:"
	@echo "  make run       - Run the application"
	@echo "  make build     - Build the application"
	@echo "  make test      - Run tests"
	@echo "  make migrate   - Run database migrations"
	@echo "  make rollback  - Rollback last migration"

run:
	go run cmd/api/main.go

build:
	go build -o bin/api cmd/api/main.go

test:
	go test -v ./...

migrate:
	goose -dir migrations postgres "postgres://postgres:password@localhost:5432/building_reports?sslmode=disable" up

rollback:
	goose -dir migrations postgres "postgres://postgres:password@localhost:5432/building_reports?sslmode=disable" down

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

install-deps:
	go mod download
	go install github.com/pressly/goose/v3/cmd/goose@latest