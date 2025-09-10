# Building Report Backend

A comprehensive backend service for building/infrastructure reporting system built with Go, following Clean Architecture, Domain-Driven Design (DDD), and Dependency Injection patterns.

## 🚀 Tech Stack

- **Go** - Programming language
- **Fiber** - Web framework
- **PostgreSQL** - Primary database
- **GORM** - ORM with raw SQL support
- **Redis** - Caching layer
- **MinIO** - Object storage for photos
- **JWT** - Authentication
- **Goose** - Database migrations

## 📋 Features

- User authentication with JWT
- CRUD operations for building reports
- Photo upload to MinIO
- Redis caching for performance
- Clean Architecture implementation
- Domain-Driven Design
- Dependency Injection
- Database migrations

## 🛠️ Setup

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- Make (optional)

### Installation

1. Clone the repository
2. Copy environment variables:
   ```bash
   cp .env.example .env