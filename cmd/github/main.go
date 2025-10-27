
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"rancher-desktop-downloads/internal/data"
)

var (
	repoOwner      = "rancher-sandbox"
	repoName       = "rancher-desktop"
	sha512sumSuffix = ".sha512sum"
	apiURL = fmt.Sprintf("https://api.github.com/repos/%s/%s/releases", repoOwner, repoName)
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

func main() {
	fmt.Println("Fetching release information...")
	releases, err := fetchReleases()
	if err != nil {
		fmt.Printf("Error fetching releases: %v\n", err)
		os.Exit(1)
	}

	if err := os.MkdirAll(data.DataDir, 0755); err != nil {
		fmt.Printf("Error creating data directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found %d releases. Processing assets...\n", len(releases))
	dailyAssetDownloads := make(map[string]int)
	for _, release := range releases {
		assetPairs := pairAssets(release.Assets)
		for assetName, downloadData := range assetPairs {
			dailyDownloads, err := data.RecordDownloadData(assetName, downloadData)
			if err != nil {
				fmt.Printf("Error recording data for %s: %v\n", assetName, err)
			}
			dailyAssetDownloads[assetName] = dailyDownloads
		}
	}

	if err := data.GenerateStatistics(dailyAssetDownloads); err != nil {
		fmt.Printf("Error generating statistics: %v\n", err)
	}

	fmt.Println("Processing complete.")
}

// fetchReleases fetches all releases from the GitHub API.
func fetchReleases() ([]Release, error) {
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
func pairAssets(assets []Asset) map[string]data.DownloadData {
	assetMap := make(map[string]Asset)
	for _, asset := range assets {
		assetMap[asset.Name] = asset
	}

	pairs := make(map[string]data.DownloadData)
	for _, asset := range assets {
		if !strings.HasSuffix(asset.Name, sha512sumSuffix) {
			checksumName := asset.Name + sha512sumSuffix
			if checksumAsset, ok := assetMap[checksumName]; ok {
				pairs[asset.Name] = data.DownloadData{
					AssetDownloads:    asset.DownloadCount,
					ChecksumDownloads: checksumAsset.DownloadCount,
				}
			}
		}
	}
	return pairs
}
