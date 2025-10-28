package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"rancher-desktop-downloads/internal/data"
)

type homebrewPackage struct {
	Name string
	Type string
}

var (
	brewAPIBaseURL = "https://formulae.brew.sh/api"
	packages       = []homebrewPackage{
		{Name: "rancher", Type: "cask"},
		{Name: "docker-desktop", Type: "cask"},
		{Name: "podman-desktop", Type: "cask"},
		{Name: "lima", Type: "formula"},
		{Name: "colima", Type: "formula"},
		{Name: "virtualbox", Type: "cask"},
		{Name: "vagrant", Type: "cask"},
		{Name: "multipass", Type: "cask"},
		{Name: "finch", Type: "cask"},
		{Name: "utm", Type: "cask"},
		{Name: "qemu", Type: "formula"},
	}
)

// Analytics represents the analytics data from the Homebrew API.
type Analytics struct {
	Analytics data.BrewAnalytics `json:"analytics"`
}

func main() {
	if err := os.MkdirAll(data.DataDir, 0755); err != nil {
		fmt.Printf("Error creating data directory: %v\n", err)
		os.Exit(1)
	}
	for _, pkg := range packages {
		fmt.Printf("Fetching Homebrew analytics data for %s (%s)...\n", pkg.Name, pkg.Type)
		analytics, err := fetchAnalyticsData(pkg)
		if err != nil {
			fmt.Printf("Error fetching analytics data for %s: %v\n", pkg.Name, err)
			continue
		}

		fmt.Printf("Processing Homebrew analytics data for %s (%s)...\n", pkg.Name, pkg.Type)
		if err := data.RecordBrewAnalytics(pkg.Name, pkg.Type, analytics.Analytics); err != nil {
			fmt.Printf("Error recording analytics for %s: %v\n", pkg.Name, err)
			continue
		}
	}

	fmt.Println("Homebrew analytics processing complete.")
}

func fetchAnalyticsData(pkg homebrewPackage) (*Analytics, error) {
	apiURL := fmt.Sprintf("%s/%s/%s.json", brewAPIBaseURL, pkg.Type, pkg.Name)
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get analytics data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	var analytics Analytics
	if err := json.NewDecoder(resp.Body).Decode(&analytics); err != nil {
		return nil, fmt.Errorf("failed to decode analytics JSON: %w", err)
	}

	return &analytics, nil
}