#!/bin/sh

# Run migrations
echo "Running database migrations..."
./server migrate

# Start the server
echo "Starting server..."
exec ./server
