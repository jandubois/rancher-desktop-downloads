
package main

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"testing"
	"time"
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
	originalDataDir := dataDir
	dataDir = tempDir
	t.Cleanup(func() { dataDir = originalDataDir })

	today := time.Now().UTC().Format("2006-01-02")
	yesterday := time.Now().UTC().AddDate(0, 0, -1).Format("2006-01-02")

	assetName := "test-asset-1.0.0.zip"
	filePath := filepath.Join(tempDir, assetName+".csv")

	t.Run("NewFile", func(t *testing.T) {
		data := DownloadData{AssetDownloads: 100, ChecksumDownloads: 10}
		if err := recordDownloadData(assetName, data); err != nil {
			t.Fatalf("recordDownloadData failed: %v", err)
		}

		expected := [][]string{
			{"date", "asset_downloads", "sha512sum_downloads"},
			{today, "100", "10"},
		}
		checkCSV(t, filePath, expected)
	})

	t.Run("ReplaceData", func(t *testing.T) {
		// This should replace the line from the previous test
		data := DownloadData{AssetDownloads: 150, ChecksumDownloads: 15}
		if err := recordDownloadData(assetName, data); err != nil {
			t.Fatalf("recordDownloadData failed: %v", err)
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
		newData := DownloadData{AssetDownloads: 200, ChecksumDownloads: 20}
		if err := recordDownloadData(assetNameForAppend, newData); err != nil {
			t.Fatalf("recordDownloadData failed: %v", err)
		}

		expected := [][]string{
			{"date", "asset_downloads", "sha512sum_downloads"},
			{yesterday, "50", "5"},
			{today, "200", "20"},
		}
		checkCSV(t, filePathForAppend, expected)
	})
}

// We need to make a small change to recordDownloadData to allow the test to override the dataDir.
// The original function is fine, but for testing this is a common pattern.
// The alternative is to pass dataDir as an argument, which is a bigger change.
// Let's stick with the global override for now as it's confined to the test.
