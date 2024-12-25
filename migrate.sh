#!/bin/bash

set -e
echo "Starting database migration..."

# Run migrations
/app/migrate -path /app/db/migration -database "$DB_SOURCE" -verbose up

echo "Migration completed successfully."