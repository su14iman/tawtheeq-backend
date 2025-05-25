package utils

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/skip2/go-qrcode"
)

func GenerateQRCodeImage(id string) ([]byte, error) {

	_ = godotenv.Load()

	baseURL := os.Getenv("FRONTEND_VERIFY_URL")
	if baseURL == "" {
		baseURL = "http://localhost:3000/verify"
	}

	fullURL := fmt.Sprintf("%s?id=%s", baseURL, id)

	var png []byte
	png, err := qrcode.Encode(fullURL, qrcode.Medium, 128)
	if err != nil {
		// return nil, fmt.Errorf("failed to generate QR code: %w", err)
		return nil, HandleError(err, "Failed to generate QR code", Error)
	}

	return png, nil
}
