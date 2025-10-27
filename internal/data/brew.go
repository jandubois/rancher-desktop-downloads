package data

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// BrewAnalytics holds the download analytics data from the Homebrew API.
type BrewAnalytics struct {
	Downloads30d  []BrewDownloadCount `json:"30d"`
	Downloads90d  []BrewDownloadCount `json:"90d"`
	Downloads365d []BrewDownloadCount `json:"365d"`
}

// BrewDownloadCount holds the download count for a specific day.
type BrewDownloadCount struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

// RecordBrewAnalytics writes the download data to the appropriate CSV file.
func RecordBrewAnalytics(pkgName, pkgType string, data BrewAnalytics) error {
	filePath := filepath.Join(DataDir, fmt.Sprintf("%s-%s.csv", pkgName, pkgType))
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
		records = [][]string{{"date", "30d", "90d", "365d"}}
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
	newRecord := []string{
		today,
		fmt.Sprintf("%d", totalDownloads(data.Downloads30d)),
		fmt.Sprintf("%d", totalDownloads(data.Downloads90d)),
		fmt.Sprintf("%d", totalDownloads(data.Downloads365d)),
	}
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

func totalDownloads(downloads []BrewDownloadCount) int {
	total := 0
	for _, dl := range downloads {
		total += dl.Count
	}
	return total
}