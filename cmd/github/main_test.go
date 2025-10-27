
package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Mock the API server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[
			{
				"name": "v1.0.0",
				"assets": [
					{
						"name": "test-asset-1.0.0.zip",
						"download_count": 100
					},
					{
						"name": "test-asset-1.0.0.zip.sha512sum",
						"download_count": 10
					}
				]
			}
		]`))
	}))
	defer server.Close()

	// Override the API URL to point to the mock server
	originalURL := apiURL
	apiURL = server.URL
	defer func() { apiURL = originalURL }()

	// Run the tests
	os.Exit(m.Run())
}

func TestFetchReleases(t *testing.T) {
	releases, err := fetchReleases()
	if err != nil {
		t.Fatalf("fetchReleases failed: %v", err)
	}

	if len(releases) != 1 {
		t.Fatalf("Expected 1 release, but got %d", len(releases))
	}

	release := releases[0]
	if release.Name != "v1.0.0" {
		t.Errorf("Expected release name to be v1.0.0, but got %s", release.Name)
	}

	if len(release.Assets) != 2 {
		t.Fatalf("Expected 2 assets, but got %d", len(release.Assets))
	}

	asset1 := release.Assets[0]
	if asset1.Name != "test-asset-1.0.0.zip" {
		t.Errorf("Expected asset name to be test-asset-1.0.0.zip, but got %s", asset1.Name)
	}
	if asset1.DownloadCount != 100 {
		t.Errorf("Expected download count to be 100, but got %d", asset1.DownloadCount)
	}

	asset2 := release.Assets[1]
	if asset2.Name != "test-asset-1.0.0.zip.sha512sum" {
		t.Errorf("Expected asset name to be test-asset-1.0.0.zip.sha512sum, but got %s", asset2.Name)
	}
	if asset2.DownloadCount != 10 {
		t.Errorf("Expected download count to be 10, but got %d", asset2.DownloadCount)
	}
}
