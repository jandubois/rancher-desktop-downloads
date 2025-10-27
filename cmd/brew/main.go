package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"rancher-desktop-downloads/internal/data"
)

var (
	brewAPIURL = "https://formulae.brew.sh/api/cask/rancher.json"
)

// Cask represents the cask data from the Homebrew API.
type Cask struct {
	Analytics data.BrewAnalytics `json:"analytics"`
}

func main() {
	fmt.Println("Fetching Homebrew analytics data...")
	cask, err := fetchCaskData()
	if err != nil {
		fmt.Printf("Error fetching cask data: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Processing Homebrew analytics data...")
	if err := os.MkdirAll(data.DataDir, 0755); err != nil {
		fmt.Printf("Error creating data directory: %v\n", err)
		os.Exit(1)
	}
	if err := data.RecordBrewAnalytics(cask.Analytics); err != nil {
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
