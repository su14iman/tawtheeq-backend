package controllers

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"tawtheeq-backend/utils"
)

func PublicKey() (*rsa.PublicKey, error) {

	keyPath := os.Getenv("PUBLIC_KEY_PATH")
	keyBytes, err := os.ReadFile(keyPath)

	if err != nil {
		return nil, utils.HandleError(err, "Failed to read public key", utils.Error)
	}
	block, _ := pem.Decode(keyBytes)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, utils.HandleError(err, "Invalid public key format", utils.Error)
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, utils.HandleError(err, "Failed to parse public key", utils.Error)
	}

	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, utils.HandleError(err, "Public key is not of type rsa.PublicKey", utils.Error)
	}
	return rsaPublicKey, nil
}

func PrivateKey() (*rsa.PrivateKey, error) {
	keyPath := os.Getenv("PRIVATE_KEY_PATH")
	keyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, utils.HandleError(err, "Failed to read private key", utils.Error)
	}

	block, _ := pem.Decode(keyBytes)
	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, utils.HandleError(err, "Invalid private key format", utils.Error)
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, utils.HandleError(err, "Failed to parse private key", utils.Error)
	}

	rsaPrivateKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, utils.HandleError(err, "Private key is not of type rsa.PrivateKey", utils.Error)
	}

	return rsaPrivateKey, nil
}
