#!/bin/bash

# Setup script for Weaviate with Docker

set -e

echo "🚀 Setting up Weaviate for alaala..."

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker first:"
    echo "   https://docs.docker.com/get-docker/"
    exit 1
fi

# Check if Weaviate container already exists
if docker ps -a --format '{{.Names}}' | grep -q "^weaviate$"; then
    echo "⚠️  Weaviate container already exists."
    read -p "Do you want to remove it and create a new one? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo "🗑️  Removing existing container..."
        docker stop weaviate 2>/dev/null || true
        docker rm weaviate 2>/dev/null || true
    else
        echo "✅ Using existing container."
        echo "   To start it: docker start weaviate"
        exit 0
    fi
fi

# Create data directory
echo "📁 Creating data directory..."
mkdir -p ~/.alaala/weaviate-data

# Start Weaviate container
echo "🐳 Starting Weaviate container..."
docker run -d \
  --name weaviate \
  -p 8080:8080 \
  -e QUERY_DEFAULTS_LIMIT=25 \
  -e AUTHENTICATION_ANONYMOUS_ACCESS_ENABLED=true \
  -e PERSISTENCE_DATA_PATH='/var/lib/weaviate' \
  -e DEFAULT_VECTORIZER_MODULE='none' \
  -e ENABLE_MODULES='' \
  -v ~/.alaala/weaviate-data:/var/lib/weaviate \
  weaviate/weaviate:latest

# Wait for Weaviate to be ready
echo "⏳ Waiting for Weaviate to be ready..."
max_attempts=30
attempt=0

while [ $attempt -lt $max_attempts ]; do
    if curl -s http://localhost:8080/v1/.well-known/ready > /dev/null 2>&1; then
        echo "✅ Weaviate is ready!"
        break
    fi
    
    attempt=$((attempt + 1))
    if [ $attempt -eq $max_attempts ]; then
        echo "❌ Weaviate failed to start after ${max_attempts} attempts"
        echo "   Check logs with: docker logs weaviate"
        exit 1
    fi
    
    echo "   Attempt $attempt/$max_attempts..."
    sleep 2
done

echo ""
echo "✨ Weaviate setup complete!"
echo ""
echo "Useful commands:"
echo "  Start Weaviate:  docker start weaviate"
echo "  Stop Weaviate:   docker stop weaviate"
echo "  View logs:       docker logs weaviate"
echo "  Remove:          docker rm -f weaviate"
echo ""
echo "Weaviate is running at: http://localhost:8080"
echo "Dashboard: http://localhost:8080/v1/meta"

