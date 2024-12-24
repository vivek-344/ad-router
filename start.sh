#!/bin/sh
set -e

if [ -z "$DB_SOURCE" ]; then
  echo "Error: DB_SOURCE is not set. Please provide the database connection string."
  exit 1
fi

echo "Running database migration..."
/app/migrate -path /app/db/migration -database "$DB_SOURCE" -verbose up

echo "Starting the application..."
exec "$@"