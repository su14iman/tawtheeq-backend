#!/bin/bash

PRIVATE_KEY_PATH="assets/keys/private.pem"
PUBLIC_KEY_PATH="assets/keys/public.pem"

# Create directory if it doesn't exist
mkdir -p assets/keys

# Generate private key if not exists
if [ ! -f "$PRIVATE_KEY_PATH" ]; then
  echo "🔐 Generating RSA private key..."
  openssl genpkey -algorithm RSA -out "$PRIVATE_KEY_PATH" -pkeyopt rsa_keygen_bits:2048
else
  echo "✅ Private key already exists"
fi

# Generate public key if not exists
if [ ! -f "$PUBLIC_KEY_PATH" ]; then
  echo "📤 Generating RSA public key..."
  openssl rsa -pubout -in "$PRIVATE_KEY_PATH" -out "$PUBLIC_KEY_PATH"
else
  echo "✅ Public key already exists"
fi

