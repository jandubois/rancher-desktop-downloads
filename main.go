
package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

var (
	repoOwner      = "rancher-sandbox"
	repoName       = "rancher-desktop"
	dataDir        = "data"
	sha512sumSuffix = ".sha512sum"
)

// Release represents a GitHub release.
type Release struct {
	Name   string  `json:"name"`
	Assets []Asset `json:"assets"`
}

// Asset represents a release asset.
type Asset struct {
	Name          string `json:"name"`
	DownloadCount int    `json:"download_count"`
}

// DownloadData holds the combined download counts for an asset and its checksum.
type DownloadData struct {
	AssetDownloads    int
	ChecksumDownloads int
}

// AssetDownloadStat holds the name and download count of an asset.
type AssetDownloadStat struct {
	Name          string
	DownloadCount int
}

func main() {
	fmt.Println("Fetching release information...")
	releases, err := fetchReleases()
	if err != nil {
		fmt.Printf("Error fetching releases: %v\n", err)
		os.Exit(1)
	}

	if err := os.MkdirAll(dataDir, 0755); err != nil {
		fmt.Printf("Error creating data directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found %d releases. Processing assets...\n", len(releases))
	allAssetDownloads := make(map[string]int)
	for _, release := range releases {
		assetPairs := pairAssets(release.Assets)
		for assetName, data := range assetPairs {
			allAssetDownloads[assetName] = data.AssetDownloads
			if err := recordDownloadData(assetName, data); err != nil {
				fmt.Printf("Error recording data for %s: %v\n", assetName, err)
			}
		}
	}

	if err := generateStatistics(allAssetDownloads); err != nil {
		fmt.Printf("Error generating statistics: %v\n", err)
	}

	fmt.Println("Processing complete.")
}

// fetchReleases fetches all releases from the GitHub API.
func fetchReleases() ([]Release, error) {
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases", repoOwner, repoName)
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get release data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("received non-200 status code %d: %s", resp.StatusCode, body)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var releases []Release
	if err := json.Unmarshal(body, &releases); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	return releases, nil
}

// pairAssets groups main assets with their corresponding sha512sum assets.
func pairAssets(assets []Asset) map[string]DownloadData {
	assetMap := make(map[string]Asset)
	for _, asset := range assets {
		assetMap[asset.Name] = asset
	}

	pairs := make(map[string]DownloadData)
	for _, asset := range assets {
		if !strings.HasSuffix(asset.Name, sha512sumSuffix) {
			checksumName := asset.Name + sha512sumSuffix
			if checksumAsset, ok := assetMap[checksumName]; ok {
				pairs[asset.Name] = DownloadData{
					AssetDownloads:    asset.DownloadCount,
					ChecksumDownloads: checksumAsset.DownloadCount,
				}
			}
		}
	}
	return pairs
}

// recordDownloadData writes the download data to the appropriate CSV file.
func recordDownloadData(assetName string, data DownloadData) error {
	filePath := filepath.Join(dataDir, assetName+".csv")
	today := time.Now().UTC().Format("2006-01-02")

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
	if len(records) == 0 {
		records = append(records, []string{"date", "asset_downloads", "sha512sum_downloads"})
	} else {
		// Check the last record
		lastRecord := records[len(records)-1]
		if lastRecord[0] == today {
			// Same day, overwrite last record
			records = records[:len(records)-1]
		}
	}

	// Append the new record
	newRecord := []string{today, fmt.Sprintf("%d", data.AssetDownloads), fmt.Sprintf("%d", data.ChecksumDownloads)}
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

// generateStatistics calculates and writes download statistics.
func generateStatistics(allAssetDownloads map[string]int) error {
	totalDownloads := 0
	assetStats := make([]AssetDownloadStat, 0, len(allAssetDownloads))

	for name, count := range allAssetDownloads {
		totalDownloads += count
		assetStats = append(assetStats, AssetDownloadStat{Name: name, DownloadCount: count})
	}

	sort.Slice(assetStats, func(i, j int) bool {
		return assetStats[i].DownloadCount > assetStats[j].DownloadCount
	})

	statsFile, err := os.Create("daily_stats.txt")
	if err != nil {
		return fmt.Errorf("failed to create statistics file: %w", err)
	}
	defer statsFile.Close()

	statsFile.WriteString(fmt.Sprintf("Total asset downloads today: %d\n\n", totalDownloads))
	statsFile.WriteString("Top 10 assets by download count:\n")
	for i := 0; i < 10 && i < len(assetStats); i++ {
		statsFile.WriteString(fmt.Sprintf("- %s: %d\n", assetStats[i].Name, assetStats[i].DownloadCount))
	}

	return nil
}
