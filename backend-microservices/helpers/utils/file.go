package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// ReadFile reads the content of a file
func ReadFile(path string) (string, error) {
	// Get the absolute path to the file
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	// Open the file
	file, err := os.Open(absPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Read the file content
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// Helper function to replace placeholders in templates
func replacePlaceholder(content, placeholder, value string) string {
	return strings.Replace(content, placeholder, value, -1)
}
