#!/bin/bash

set -e
echo "Starting database migration..."

# Load env and run migrations
source /app/app.env
/app/migrate -path /app/db/migration -database "$DB_SOURCE" -verbose up

echo "Migration completed successfully."