package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

var (
	logFile     *os.File
	logFilePath = "logs/app.log"
	logLimitMB  = 100
)

func InitLogging() error {
	if os.Getenv("LOGGING_ENABLE") != "true" {
		return nil
	}

	logDir := filepath.Dir(logFilePath)
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	if limit := os.Getenv("LOGGING_LIMIT"); limit != "" {
		if val, err := strconv.Atoi(limit); err == nil {
			logLimitMB = val
		}
	}

	// Rotate if needed
	if stat, err := os.Stat(logFilePath); err == nil {
		if stat.Size() >= int64(logLimitMB)*1024*1024 {
			backup := filepath.Join(logDir, fmt.Sprintf("app-%s.log", stat.ModTime().Format("20060102-150405")))
			_ = os.Rename(logFilePath, backup)
		}
	}

	f, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	logFile = f
	log.SetOutput(logFile)

	return nil
}

func CloseLogging() {
	if logFile != nil {
		logFile.Close()
	}
}
