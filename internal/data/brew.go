
package data

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"
)

var (
	// BrewCSV is the path to the brew analytics csv file.
	BrewCSV = "data/homebrew_analytics.csv"
)

// BrewAnalytics represents the analytics data from the Homebrew API.
type BrewAnalytics struct {
	Install struct {
		ThirtyDays      map[string]int `json:"30d"`
		NinetyDays      map[string]int `json:"90d"`
		ThreeHundredSixtyFiveDays map[string]int `json:"365d"`
	} `json:"install"`
}

// RecordBrewAnalytics writes the brew analytics data to the CSV file.
func RecordBrewAnalytics(analytics BrewAnalytics) error {
	today := time.Now().Format("2006-01-02")
	thirtyDayCount := analytics.Install.ThirtyDays["rancher"]
	ninetyDayCount := analytics.Install.NinetyDays["rancher"]
	threeHundredSixtyFiveDayCount := analytics.Install.ThreeHundredSixtyFiveDays["rancher"]

	file, err := os.OpenFile(BrewCSV, os.O_RDWR|os.O_CREATE, 0644)
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
	if len(records) == 0 {
		records = [][]string{{"date", "30d", "90d", "365d"}}
	} else {
		// Check the last record
		lastRecord := records[len(records)-1]
		if lastRecord[0] == today {
			// Same day, overwrite last record
			records = records[:len(records)-1]
		}
	}

	// Append the new record
	newRecord := []string{
		today,
		fmt.Sprintf("%d", thirtyDayCount),
		fmt.Sprintf("%d", ninetyDayCount),
		fmt.Sprintf("%d", threeHundredSixtyFiveDayCount),
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
