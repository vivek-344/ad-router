#!/bin/bash

set -e
echo "Starting api service..."

# Load env variables
source /app/app.env

echo "Api service running..."
exec "$@"
