#!/bin/bash

set -e

echo "🚀 Starting Tawtheeq Backend..."

# Create needed directories if they don't exist
mkdir -p docker-data/mysql
mkdir -p docker-data/redis
mkdir -p docker-data/minio
mkdir -p uploads
mkdir -p logs

# Check if .env file exists
if [ ! -f .env ]; then
  echo "❌ .env file not found. Please create it first."
  exit 1
fi

# Choose the correct docker compose command
if command -v docker-compose &> /dev/null; then
  COMPOSE_COMMAND="docker-compose"
elif docker compose version &> /dev/null; then
  COMPOSE_COMMAND="docker compose"
else
  echo "❌ Neither docker-compose nor docker compose is installed."
  exit 1
fi

# Parse parameters
BUILD_ARG=""
DETACHED_ARG=""

for arg in "$@"; do
  case $arg in
    build)
      BUILD_ARG="--build"
      ;;
    -d|--detached)
      DETACHED_ARG="-d"
      ;;
  esac
done

echo "⚡ Running: $COMPOSE_COMMAND up $BUILD_ARG $DETACHED_ARG"
$COMPOSE_COMMAND up $BUILD_ARG $DETACHED_ARG
