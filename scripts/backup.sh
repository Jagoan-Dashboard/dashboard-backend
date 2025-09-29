#!/bin/bash

# Database Backup Script
set -e

BACKUP_DIR="/opt/building-report-backend/backups"
DATE=$(date +%Y%m%d_%H%M%S)
DB_CONTAINER="building-report-db"

# Create backup directory if it doesn't exist
mkdir -p "$BACKUP_DIR"

echo "🗄️ Starting database backup..."

# Check if database container is running
if ! docker ps | grep -q "$DB_CONTAINER"; then
    echo "❌ Database container is not running"
    exit 1
fi

# Create database backup
BACKUP_FILE="$BACKUP_DIR/backup_${DATE}.sql"
echo "📦 Creating backup: $BACKUP_FILE"

docker-compose exec -T postgres pg_dump -U postgres building_reports > "$BACKUP_FILE"

if [ $? -eq 0 ]; then
    echo "✅ Database backup created successfully: $BACKUP_FILE"

    # Compress the backup
    gzip "$BACKUP_FILE"
    echo "🗜️ Backup compressed: ${BACKUP_FILE}.gz"

    # Keep only last 7 backups
    cd "$BACKUP_DIR"
    ls -t backup_*.sql.gz | tail -n +8 | xargs -r rm --
    echo "🧹 Old backups cleaned up (keeping last 7)"

    # Show backup size
    BACKUP_SIZE=$(du -h "${BACKUP_FILE}.gz" | cut -f1)
    echo "💾 Backup size: $BACKUP_SIZE"

else
    echo "❌ Database backup failed"
    exit 1
fi

echo "🎉 Backup completed successfully!"