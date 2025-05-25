package utils

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

const (
	Info    = "INFO"
	Warning = "WARNING"
	Error   = "ERROR"
)

func HandleError(err error, customMessage string, level string) error {
	if err == nil {
		return nil
	}

	if os.Getenv("LOGGING_ENABLE") == "true" && shouldLogLevel(level) {
		logMessage := fmt.Sprintf(
			"\n%s [%s] : %s : %v\n",
			time.Now().Format(time.RFC3339),
			level,
			customMessage,
			err,
		)
		log.Println(logMessage)
	}

	return fmt.Errorf("%s: %w", customMessage, err)
}

func shouldLogLevel(level string) bool {
	envLevel := strings.ToUpper(os.Getenv("LOGGING_LEVEL"))

	if envLevel == "*" {
		return true
	}

	levelPriority := map[string]int{
		"INFO":    1,
		"WARNING": 2,
		"ERROR":   3,
	}

	current := levelPriority[strings.ToUpper(level)]
	threshold := levelPriority[envLevel]

	return current >= threshold
}

func MainErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}
	return c.Status(code).JSON(fiber.Map{
		"error":  err.Error(),
		"status": code,
	})
}
