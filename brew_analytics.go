package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	brewAPIURL = "https://formulae.brew.sh/api/cask/rancher.json"
	brewCSV    = "data/homebrew_analytics.csv"
)

// BrewAnalytics represents the analytics data from the Homebrew API.
type BrewAnalytics struct {
	Install struct {
		ThirtyDays      map[string]int `json:"30d"`
		NinetyDays      map[string]int `json:"90d"`
		ThreeHundredSixtyFiveDays map[string]int `json:"365d"`
	} `json:"install"`
}

// Cask represents the cask data from the Homebrew API.
type Cask struct {
	Analytics BrewAnalytics `json:"analytics"`
}

func main() {
	fmt.Println("Fetching Homebrew analytics data...")
	cask, err := fetchCaskData()
	if err != nil {
		fmt.Printf("Error fetching cask data: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Processing Homebrew analytics data...")
	if err := recordAnalytics(cask.Analytics); err != nil {
		fmt.Printf("Error recording analytics: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Homebrew analytics processing complete.")
}

func fetchCaskData() (*Cask, error) {
	resp, err := http.Get(brewAPIURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get cask data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	var cask Cask
	if err := json.NewDecoder(resp.Body).Decode(&cask); err != nil {
		return nil, fmt.Errorf("failed to decode cask JSON: %w", err)
	}

	return &cask, nil
}

func recordAnalytics(analytics BrewAnalytics) error {
	today := time.Now().Format("2006-01-02")
	thirtyDayCount := analytics.Install.ThirtyDays["rancher"]
	ninetyDayCount := analytics.Install.NinetyDays["rancher"]
	threeHundredSixtyFiveDayCount := analytics.Install.ThreeHundredSixtyFiveDays["rancher"]

	file, err := os.OpenFile(brewCSV, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Check if file is new, write header if so
	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}
	if info.Size() == 0 {
		if err := writer.Write([]string{"date", "30d", "90d", "365d"}); err != nil {
			return fmt.Errorf("failed to write header: %w", err)
		}
	}

	record := []string{
		today,
		fmt.Sprintf("%d", thirtyDayCount),
		fmt.Sprintf("%d", ninetyDayCount),
		fmt.Sprintf("%d", threeHundredSixtyFiveDayCount),
	}

	if err := writer.Write(record); err != nil {
		return fmt.Errorf("failed to write record: %w", err)
	}

	return nil
}
