
package data

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

var (
	// DataDir is the directory where the data files are stored.
	DataDir   = "data"
	// PacificTZ is the timezone used for the data files.
	PacificTZ *time.Location
)

func init() {
	var err error
	PacificTZ, err = time.LoadLocation("America/Los_Angeles")
	if err != nil {
		fmt.Printf("Error loading location: %v\n", err)
		os.Exit(1)
	}
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

// RecordDownloadData writes the download data to the appropriate CSV file and returns the daily download count.
func RecordDownloadData(assetName string, data DownloadData) (int, error) {
	filePath := filepath.Join(DataDir, assetName+".csv")
	today := time.Now().In(PacificTZ).Format("2006-01-02")
	dailyDownloads := 0

	// Read existing data
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return 0, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return 0, fmt.Errorf("failed to read csv records: %w", err)
	}

	// Check if file is new, write header if so
	if len(records) <= 1 {
		records = [][]string{{"date", "asset_downloads", "sha512sum_downloads"}}
		dailyDownloads = data.AssetDownloads
	} else {
		// Check the last record
		lastRecord := records[len(records)-1]
		if lastRecord[0] == today {
			// Same day, overwrite last record
			records = records[:len(records)-1]
			if len(records) > 1 {
				prevRecord := records[len(records)-1]
				prevDownloads, _ := strconv.Atoi(prevRecord[1])
				dailyDownloads = data.AssetDownloads - prevDownloads
			} else {
				dailyDownloads = data.AssetDownloads
			}
		} else {
			prevDownloads, _ := strconv.Atoi(lastRecord[1])
			dailyDownloads = data.AssetDownloads - prevDownloads
		}
	}

	// Append the new record
	newRecord := []string{today, fmt.Sprintf("%d", data.AssetDownloads), fmt.Sprintf("%d", data.ChecksumDownloads)}
	records = append(records, newRecord)

	// Write all records back to the file
	if err := file.Truncate(0); err != nil {
		return 0, fmt.Errorf("failed to truncate file: %w", err)
	}
	if _, err := file.Seek(0, 0); err != nil {
		return 0, fmt.Errorf("failed to seek file: %w", err)
	}

	writer := csv.NewWriter(file)
	if err := writer.WriteAll(records); err != nil {
		return 0, fmt.Errorf("failed to write csv records: %w", err)
	}

	return dailyDownloads, nil
}

// GenerateStatistics calculates and writes download statistics.
func GenerateStatistics(dailyAssetDownloads map[string]int) error {
	totalDailyDownloads := 0
	assetStats := make([]AssetDownloadStat, 0, len(dailyAssetDownloads))

	for name, count := range dailyAssetDownloads {
		totalDailyDownloads += count
		assetStats = append(assetStats, AssetDownloadStat{Name: name, DownloadCount: count})
	}

	sort.Slice(assetStats, func(i, j int) bool {
		if assetStats[i].DownloadCount == assetStats[j].DownloadCount {
			return assetStats[i].Name < assetStats[j].Name
		}
		return assetStats[i].DownloadCount > assetStats[j].DownloadCount
	})

	statsFile, err := os.Create("daily_stats.txt")
	if err != nil {
		return fmt.Errorf("failed to create statistics file: %w", err)
	}
	defer statsFile.Close()

	statsFile.WriteString(fmt.Sprintf("Total asset downloads today: %d\n\n", totalDailyDownloads))
	statsFile.WriteString("Top 10 assets by daily download count:\n")
	for i := 0; i < 10 && i < len(assetStats); i++ {
		statsFile.WriteString(fmt.Sprintf("- %s: %d\n", assetStats[i].Name, assetStats[i].DownloadCount))
	}

	return nil
}
