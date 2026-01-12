package utils

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// GenerateCSV creates a CSV file from MongoDB records (bson.M)
func GenerateCSV(records []bson.M, collectionName string, outputDir string) (string, error) {
	if len(records) == 0 {
		// Create an empty CSV file even if there are no records
		fileName := fmt.Sprintf("%s_%s.csv", collectionName, time.Now().Format("2006-01-02"))
		filePath := filepath.Join(outputDir, fileName)
		
		file, err := os.Create(filePath)
		if err != nil {
			return "", fmt.Errorf("failed to create CSV file: %w", err)
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		// Write empty header row
		if err := writer.Write([]string{"No data available"}); err != nil {
			return "", fmt.Errorf("failed to write CSV header: %w", err)
		}

		return filePath, nil
	}

	// Get all unique keys from all records
	keySet := make(map[string]bool)
	for _, record := range records {
		for key := range record {
			keySet[key] = true
		}
	}

	// Sort keys for consistent column order
	keys := make([]string, 0, len(keySet))
	for key := range keySet {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Create CSV file
	fileName := fmt.Sprintf("%s_%s.csv", collectionName, time.Now().Format("2006-01-02"))
	filePath := filepath.Join(outputDir, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	if err := writer.Write(keys); err != nil {
		return "", fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data rows
	for _, record := range records {
		row := make([]string, len(keys))
		for i, key := range keys {
			value := record[key]
			row[i] = convertValueToString(value)
		}
		if err := writer.Write(row); err != nil {
			return "", fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	return filePath, nil
}

// convertValueToString converts various MongoDB BSON types to string
func convertValueToString(value interface{}) string {
	if value == nil {
		return ""
	}

	switch v := value.(type) {
	case string:
		return v
	case int:
		return fmt.Sprintf("%d", v)
	case int32:
		return fmt.Sprintf("%d", v)
	case int64:
		return fmt.Sprintf("%d", v)
	case float32:
		return fmt.Sprintf("%f", v)
	case float64:
		return fmt.Sprintf("%f", v)
	case bool:
		return fmt.Sprintf("%t", v)
	case time.Time:
		return v.Format(time.RFC3339Nano)
	case bson.M:
		// Convert nested document to JSON string
		return fmt.Sprintf("%v", v)
	case []interface{}:
		// Convert array to JSON-like string
		return fmt.Sprintf("%v", v)
	case bson.A:
		return fmt.Sprintf("%v", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// EnsureDir creates a directory if it doesn't exist
func EnsureDir(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return os.MkdirAll(dirPath, 0755)
	}
	return nil
}

