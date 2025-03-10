package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

func ToCache(path string, data interface{}) error {
	// Create or truncate the file
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Setup the encoder based on whether indentation is requested
	enc := json.NewEncoder(file)

	// Encode and write the data to the file
	if err := enc.Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}

func FromCache(path string, target interface{}) (float64, error) {
	// Get file info to check modification time
	fileInfo, err := os.Stat(path)
	if err != nil {
		return 0, fmt.Errorf("failed to get file info: %w", err)
	}

	// Calculate minutes since last modification
	modTime := fileInfo.ModTime()
	minutesAgo := time.Since(modTime).Minutes()

	// Open the file
	file, err := os.Open(path)
	if err != nil {
		return minutesAgo, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Create a decoder and decode into the target
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(target); err != nil {
		return minutesAgo, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return minutesAgo, nil
}
