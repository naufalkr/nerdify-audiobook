package utils

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

// InitLogger sets up logging to write to both stdout and a log file
func InitLogger() {
	// Create logs directory if it doesn't exist
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Printf("Failed to create log directory: %v", err)
		return
	}

	// Open log file
	logFile, err := os.OpenFile(
		filepath.Join(logDir, "app.log"),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		log.Printf("Failed to open log file: %v", err)
		return
	}
	// Set log output to both file and stdout
	log.SetOutput(os.Stdout)
	if logFile != nil {
		multiWriter := io.MultiWriter(os.Stdout, logFile)
		log.SetOutput(multiWriter)
	}
}
