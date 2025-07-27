#!/bin/bash

# Age Verification Demo Server Startup Script
echo "🔐 Starting Age Verification Demo Server..."

# Find available port
PORT=8089
while lsof -ti:$PORT > /dev/null 2>&1; do
    PORT=$((PORT + 1))
    if [ $PORT -gt 8099 ]; then
        echo "❌ No available ports found in range 8089-8099"
        exit 1
    fi
done

echo "🚀 Starting server on port $PORT..."
echo "🌐 Age Verification Demo: http://localhost:$PORT/age-verification.html"
echo "📊 Main Demo: http://localhost:$PORT/"
echo "🏥 Health Check: http://localhost:$PORT/health"
echo ""
echo "Press Ctrl+C to stop the server"

# Start the server
go run ./cmd/server --port=$PORT
