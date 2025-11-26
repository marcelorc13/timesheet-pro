#!/bin/sh
set -e

echo "Starting migration process..."

# Check if POSTGRES_URL is set
if [ -z "$POSTGRES_URL" ]; then
    echo "ERROR: POSTGRES_URL environment variable is not set"
    exit 1
fi

# Run database migrations
echo "Running database migrations..."
cd /app/migrations
goose postgres "$POSTGRES_URL" up

if [ $? -eq 0 ]; then
    echo "Migrations completed successfully!"
else
    echo "ERROR: Migrations failed!"
    exit 1
fi

# Go back to root directory
cd /root

# Start the application
echo "Starting application..."
exec ./main
