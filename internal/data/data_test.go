
package data_test

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"testing"
	"time"

	"rancher-desktop-downloads/internal/data"
)

// Helper function to check CSV content
func checkCSV(t *testing.T, path string, expected [][]string) {
	t.Helper()
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("Failed to open file %s: %v", path, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("Failed to read CSV from %s: %v", path, err)
	}

	if len(records) != len(expected) {
		t.Fatalf("Expected %d records, but got %d", len(expected), len(records))
	}

	for i, expectedRecord := range expected {
		actualRecord := records[i]
		if len(actualRecord) != len(expectedRecord) {
			t.Errorf("Record #%d: expected %d fields, got %d", i, len(expectedRecord), len(actualRecord))
			continue
		}
		for j, expectedField := range expectedRecord {
			if actualRecord[j] != expectedField {
				t.Errorf("Record #%d, Field #%d: expected %s, got %s", i, j, expectedField, actualRecord[j])
			}
		}
	}
}

func TestRecordDownloadData(t *testing.T) {
	tempDir := t.TempDir()
	// Override the dataDir for testing purposes
	originalDataDir := data.DataDir
	data.DataDir = tempDir
	t.Cleanup(func() { data.DataDir = originalDataDir })

	today := time.Now().In(data.RecordTZ).Format("2006-01-02")
	yesterday := time.Now().In(data.RecordTZ).AddDate(0, 0, -1).Format("2006-01-02")

	assetName := "test-asset-1.0.0.zip"
	filePath := filepath.Join(tempDir, assetName+".csv")

	t.Run("NewFile", func(t *testing.T) {
		downloadData := data.DownloadData{AssetDownloads: 100, ChecksumDownloads: 10}
		if _, err := data.RecordDownloadData(assetName, downloadData); err != nil {
			t.Fatalf("RecordDownloadData failed: %v", err)
		}

		expected := [][]string{
			{"date", "asset_downloads", "sha512sum_downloads"},
			{today, "100", "10"},
		}
		checkCSV(t, filePath, expected)
	})

	t.Run("ReplaceData", func(t *testing.T) {
		// This should replace the line from the previous test
		downloadData := data.DownloadData{AssetDownloads: 150, ChecksumDownloads: 15}
		if _, err := data.RecordDownloadData(assetName, downloadData); err != nil {
			t.Fatalf("RecordDownloadData failed: %v", err)
		}

		expected := [][]string{
			{"date", "asset_downloads", "sha512sum_downloads"},
			{today, "150", "15"},
		}
		checkCSV(t, filePath, expected)
	})

	t.Run("AppendData", func(t *testing.T) {
		// To test appending, we need to manually create a file with an older date
		assetNameForAppend := "test-asset-2.0.0.zip"
		filePathForAppend := filepath.Join(tempDir, assetNameForAppend+".csv")

		// Create a file with yesterday's data
		initialRecords := [][]string{
			{"date", "asset_downloads", "sha512sum_downloads"},
			{yesterday, "50", "5"},
		}
		file, _ := os.Create(filePathForAppend)
		writer := csv.NewWriter(file)
		writer.WriteAll(initialRecords)
		writer.Flush()
		file.Close()

		// Now, record today's data
		newData := data.DownloadData{AssetDownloads: 200, ChecksumDownloads: 20}
		if _, err := data.RecordDownloadData(assetNameForAppend, newData); err != nil {
			t.Fatalf("RecordDownloadData failed: %v", err)
		}

		expected := [][]string{
			{"date", "asset_downloads", "sha512sum_downloads"},
			{yesterday, "50", "5"},
			{today, "200", "20"},
		}
		checkCSV(t, filePathForAppend, expected)
	})
}

func TestRecordBrewAnalytics(t *testing.T) {
	tempDir := t.TempDir()
	// Override the BrewCSV for testing purposes
	originalBrewCSV := data.BrewCSV
	data.BrewCSV = filepath.Join(tempDir, "homebrew_analytics.csv")
	t.Cleanup(func() { data.BrewCSV = originalBrewCSV })
	today := time.Now().In(data.RecordTZ).Format("2006-01-02")
	yesterday := time.Now().In(data.RecordTZ).AddDate(0, 0, -1).Format("2006-01-02")
	filePath := filepath.Join(tempDir, "homebrew_analytics.csv")

	t.Run("NewFile", func(t *testing.T) {
		analytics := data.BrewAnalytics{}
		analytics.Install.ThirtyDays = map[string]int{"rancher": 100}
		analytics.Install.NinetyDays = map[string]int{"rancher": 200}
		analytics.Install.ThreeHundredSixtyFiveDays = map[string]int{"rancher": 300}

		if err := data.RecordBrewAnalytics(analytics); err != nil {
			t.Fatalf("RecordBrewAnalytics failed: %v", err)
		}

		expected := [][]string{
			{"date", "30d", "90d", "365d"},
			{today, "100", "200", "300"},
		}
		checkCSV(t, filePath, expected)
	})

	t.Run("ReplaceData", func(t *testing.T) {
		// This should replace the line from the previous test
		analytics := data.BrewAnalytics{}
		analytics.Install.ThirtyDays = map[string]int{"rancher": 150}
		analytics.Install.NinetyDays = map[string]int{"rancher": 250}
		analytics.Install.ThreeHundredSixtyFiveDays = map[string]int{"rancher": 350}

		if err := data.RecordBrewAnalytics(analytics); err != nil {
			t.Fatalf("RecordBrewAnalytics failed: %v", err)
		}

		expected := [][]string{
			{"date", "30d", "90d", "365d"},
			{today, "150", "250", "350"},
		}
		checkCSV(t, filePath, expected)
	})

	t.Run("AppendData", func(t *testing.T) {
		// To test appending, we need to manually create a file with an older date
		initialRecords := [][]string{
			{"date", "30d", "90d", "365d"},
			{yesterday, "50", "100", "150"},
		}
		file, _ := os.Create(filePath)
		writer := csv.NewWriter(file)
		writer.WriteAll(initialRecords)
		writer.Flush()
		file.Close()

		// Now, record today's data
		analytics := data.BrewAnalytics{}
		analytics.Install.ThirtyDays = map[string]int{"rancher": 200}
		analytics.Install.NinetyDays = map[string]int{"rancher": 300}
		analytics.Install.ThreeHundredSixtyFiveDays = map[string]int{"rancher": 400}

		if err := data.RecordBrewAnalytics(analytics); err != nil {
			t.Fatalf("RecordBrewAnalytics failed: %v", err)
		}

		expected := [][]string{
			{"date", "30d", "90d", "365d"},
			{yesterday, "50", "100", "150"},
			{today, "200", "300", "400"},
		}
		checkCSV(t, filePath, expected)
	})
}


// We need to make a small change to recordDownloadData to allow the test to override the dataDir.
// The original function is fine, but for testing this is a common pattern.
// The alternative is to pass dataDir as an argument, which is a bigger change.
// Let's stick with the global override for now as it's confined to the test.

func TestGenerateStatistics(t *testing.T) {
	assetDownloads := map[string]int{
		"asset-1": 100,
		"asset-2": 200,
		"asset-3": 50,
		"asset-4": 300,
		"asset-5": 150,
		"asset-6": 250,
		"asset-7": 500,
		"asset-8": 400,
		"asset-9": 350,
		"asset-10": 450,
		"asset-11": 50,
	}

	err := data.GenerateStatistics(assetDownloads)
	if err != nil {
		t.Fatalf("generateStatistics failed: %v", err)
	}

	// Read the generated file
	content, err := os.ReadFile("daily_stats.txt")
	if err != nil {
		t.Fatalf("Failed to read daily_stats.txt: %v", err)
	}
	defer os.Remove("daily_stats.txt")

			expectedContent := "Total asset downloads today: 2800\n\n" +
			"Top 10 assets by daily download count:\n" +
			"- asset-7: 500\n" +
			"- asset-10: 450\n" +
			"- asset-8: 400\n" +
			"- asset-9: 350\n" +
			"- asset-4: 300\n" +
			"- asset-6: 250\n" +
			"- asset-2: 200\n" +
			"- asset-5: 150\n" +
			"- asset-1: 100\n" +
			"- asset-11: 50\n"
		if string(content) != expectedContent {
		t.Errorf("Unexpected content in daily_stats.txt.\nGot:\n%s\nExpected:\n%s", string(content), expectedContent)
	}
}
