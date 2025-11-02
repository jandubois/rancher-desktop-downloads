package data

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const chocolateySubdir = "chocolatey"

// RecordChocolateyTotal writes the total download count to a CSV file.
func RecordChocolateyTotal(packageName string, totalDownloads int) error {
	chocolateyDir := filepath.Join(DataDir, chocolateySubdir)
	if err := os.MkdirAll(chocolateyDir, 0755); err != nil {
		return fmt.Errorf("failed to create chocolatey directory: %w", err)
	}

	filePath := filepath.Join(chocolateyDir, packageName+"-total.csv")
	today := time.Now().In(RecordTZ).Format("2006-01-02")

	// Read existing data
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read csv records: %w", err)
	}

	// Check if file is new, write header if so
	if len(records) <= 1 {
		records = [][]string{{"date", "total_downloads"}}
	}

	// Check the last record
	if len(records) > 1 {
		lastRecord := records[len(records)-1]
		if lastRecord[0] == today {
			// Same day, overwrite last record
			records = records[:len(records)-1]
		}
	}

	// Append the new record
	newRecord := []string{today, fmt.Sprintf("%d", totalDownloads)}
	records = append(records, newRecord)

	// Write all records back to the file
	if err := file.Truncate(0); err != nil {
		return fmt.Errorf("failed to truncate file: %w", err)
	}
	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to seek file: %w", err)
	}

	writer := csv.NewWriter(file)
	if err := writer.WriteAll(records); err != nil {
		return fmt.Errorf("failed to write csv records: %w", err)
	}

	return nil
}

// RecordChocolateyVersion writes the version-specific download count to a CSV file.
func RecordChocolateyVersion(packageName, version string, versionDownloads int) error {
	chocolateyDir := filepath.Join(DataDir, chocolateySubdir)
	if err := os.MkdirAll(chocolateyDir, 0755); err != nil {
		return fmt.Errorf("failed to create chocolatey directory: %w", err)
	}

	filePath := filepath.Join(chocolateyDir, fmt.Sprintf("%s-%s.csv", packageName, version))
	today := time.Now().In(RecordTZ).Format("2006-01-02")

	// Read existing data
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read csv records: %w", err)
	}

	// Check if file is new, write header if so
	if len(records) <= 1 {
		records = [][]string{{"date", "version_downloads"}}
	}

	// Check the last record
	if len(records) > 1 {
		lastRecord := records[len(records)-1]
		if lastRecord[0] == today {
			// Same day, overwrite last record
			records = records[:len(records)-1]
		}
	}

	// Append the new record
	newRecord := []string{today, fmt.Sprintf("%d", versionDownloads)}
	records = append(records, newRecord)

	// Write all records back to the file
	if err := file.Truncate(0); err != nil {
		return fmt.Errorf("failed to truncate file: %w", err)
	}
	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to seek file: %w", err)
	}

	writer := csv.NewWriter(file)
	if err := writer.WriteAll(records); err != nil {
		return fmt.Errorf("failed to write csv records: %w", err)
	}

	return nil
}
