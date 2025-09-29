#!/bin/bash

# Building Report Backend Deployment Script
set -e

echo "🚀 Starting deployment..."

# Check if running as root
if [[ $EUID -eq 0 ]]; then
   echo "❌ This script should not be run as root"
   exit 1
fi

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker first."
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose is not installed. Please install Docker Compose first."
    exit 1
fi

# Check if .env file exists
if [ ! -f ".env" ]; then
    echo "❌ .env file not found. Please copy .env.production to .env and configure it."
    exit 1
fi

# Pull latest changes
echo "📥 Pulling latest changes..."
git pull origin main

# Build and deploy
echo "🔨 Building and starting services..."
docker-compose down
docker-compose up -d --build

# Wait for services to be ready
echo "⏳ Waiting for services to start..."
sleep 10

# Check if services are running
echo "🔍 Checking service status..."
if docker-compose ps | grep -q "Up"; then
    echo "✅ Services are running!"
else
    echo "❌ Some services failed to start. Check logs:"
    docker-compose logs --tail=20
    exit 1
fi

# Health check
echo "🏥 Running health check..."
if curl -f -s http://localhost:8080/health > /dev/null; then
    echo "✅ Application is healthy!"
else
    echo "⚠️ Health check failed. Check application logs:"
    docker-compose logs app --tail=20
fi

echo "🎉 Deployment completed!"
echo ""
echo "📋 Quick commands:"
echo "  - View logs: docker-compose logs -f"
echo "  - Check status: docker-compose ps"
echo "  - Stop services: docker-compose down"
echo "  - Restart: docker-compose restart"