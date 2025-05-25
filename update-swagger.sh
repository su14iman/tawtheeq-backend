#!/bin/bash

# Make sure you have swaggo/swag installed
# If not installed, install it with:
# go install github.com/swaggo/swag/cmd/swag@latest

echo "Updating Swagger documentation..."

# Remove old docs
rm -rf ./docs/docs.go ./docs/swagger.json ./docs/swagger.yaml

# Generate new docs
swag init -g main.go --parseDependency --output ./docs

# Verify files were generated
if [ -f "./docs/docs.go" ] && [ -f "./docs/swagger.json" ] && [ -f "./docs/swagger.yaml" ]; then
    echo "✅ Swagger documentation updated successfully!"
    echo "📚 Access Swagger UI at: http://localhost:4000/swagger/index.html when server is running"
else
    echo "❌ Failed to update Swagger documentation."
    echo "Make sure swag is installed. Run: go install github.com/swaggo/swag/cmd/swag@latest"
fi
