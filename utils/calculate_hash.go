package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

// CalculateFileHash calculates the SHA256 hash of a file at the given path
func CalculateFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", HandleError(err, "Failed to open file for hashing", Error)
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", HandleError(err, "Failed to hash file content", Error)
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}
