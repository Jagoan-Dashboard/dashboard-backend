# Makefile
.PHONY: help run build test migrate
MIGRATIONS_DIR := ./migrations

# Load .env jika ada
ifneq (,$(wildcard .env))
include .env
export
endif

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

# Database Seeding
.PHONY: seed seed-users seed-reports seed-spatial seed-water seed-bina seed-agriculture seed-all seed-fresh

# Seed specific table
seed-users:
	@echo "Seeding users table..."
	@go run cmd/seeder/main.go users

seed-reports:
	@echo "Seeding reports table..."
	@go run cmd/seeder/main.go reports

seed-spatial:
	@echo "Seeding spatial planning table..."
	@go run cmd/seeder/main.go spatial_planning

seed-water:
	@echo "Seeding water resources table..."
	@go run cmd/seeder/main.go water_resources

seed-bina:
	@echo "Seeding bina marga table..."
	@go run cmd/seeder/main.go bina_marga

seed-agriculture:
	@echo "Seeding agriculture table..."
	@go run cmd/seeder/main.go agriculture

# Seed all tables
seed-all:
	@echo "Seeding all tables..."
	@go run cmd/seeder/main.go all

# Fresh seed: drop all tables, migrate, then seed
seed-fresh:
	@echo "Fresh seeding: dropping tables, migrating, and seeding..."
	@make db-reset
	@make migrate-up
	@make seed-all

# Alias for seed-all
seed: seed-all

# Database management
db-reset:
	@echo "Resetting database..."
	@PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d postgres -c "DROP DATABASE IF EXISTS $(DB_NAME);"
	@PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d postgres -c "CREATE DATABASE $(DB_NAME);"
	@echo "Database reset complete"

# Migration commands (jika belum ada)
migrate-up:
	@echo "Running migrations..."
	@goose -dir migrations postgres "host=$(DB_HOST) port=$(DB_PORT) user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(DB_NAME) sslmode=disable" up

migrate-down:
	@echo "Rolling back migrations..."
	@goose -dir migrations postgres "host=$(DB_HOST) port=$(DB_PORT) user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(DB_NAME) sslmode=disable" down

migrate-status:
	@echo "Migration status..."
	@goose -dir migrations postgres "host=$(DB_HOST) port=$(DB_PORT) user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(DB_NAME) sslmode=disable" status

migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "❌ Usage: make migrate-create name=<migration_name>"; \
	else \
		goose create $(name) sql -dir $(MIGRATIONS_DIR); \
	fi