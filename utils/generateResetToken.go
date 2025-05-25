package utils

import (
	"github.com/google/uuid"
)

func GenerateResetToken() string {
	return uuid.New().String()
}
