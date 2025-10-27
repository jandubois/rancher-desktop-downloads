
package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"rancher-desktop-downloads/internal/data"
	"strings"
	"testing"
)

var ( mockServer *httptest.Server)

func TestMain(m *testing.M) {
	mockServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if strings.HasPrefix(r.URL.Path, "/api/cask") {
			w.Write([]byte(`{
				"analytics": {
					"30d": [{"date": "2025-10-26", "count": 1234}],
					"90d": [{"date": "2025-10-26", "count": 5678}],
					"365d": [{"date": "2025-10-26", "count": 91011}]
				}
			}`))
		} else if strings.HasPrefix(r.URL.Path, "/api/formula") {
			w.Write([]byte(`{
				"analytics": {
					"30d": [{"date": "2025-10-26", "count": 4321}],
					"90d": [{"date": "2025-10-26", "count": 8765}],
					"365d": [{"date": "2025-10-26", "count": 11019}]
				}
			}`))
		}
	}))
	originalURL := brewAPIBaseURL
	brewAPIBaseURL = mockServer.URL + "/api"

	code := m.Run()

	brewAPIBaseURL = originalURL
	mockServer.Close()

	os.Exit(code)
}

func TestFetchAnalyticsData(t *testing.T) {
	caskPkg := homebrewPackage{Name: "rancher", Type: "cask"}
	caskAnalytics, err := fetchAnalyticsData(caskPkg)
	if err != nil {
		t.Fatalf("fetchAnalyticsData for cask failed: %v", err)
	}
	if totalDownloads(caskAnalytics.Analytics.Downloads30d) != 1234 {
		t.Errorf("Expected 30d count for cask to be 1234, but got %d", totalDownloads(caskAnalytics.Analytics.Downloads30d))
	}

	formulaPkg := homebrewPackage{Name: "lima", Type: "formula"}
	formulaAnalytics, err := fetchAnalyticsData(formulaPkg)
	if err != nil {
		t.Fatalf("fetchAnalyticsData for formula failed: %v", err)
	}
	if totalDownloads(formulaAnalytics.Analytics.Downloads30d) != 4321 {
		t.Errorf("Expected 30d count for formula to be 4321, but got %d", totalDownloads(formulaAnalytics.Analytics.Downloads30d))
	}
}

func totalDownloads(downloads []data.BrewDownloadCount) int {
	total := 0
	for _, dl := range downloads {
		total += dl.Count
	}
	return total
}
