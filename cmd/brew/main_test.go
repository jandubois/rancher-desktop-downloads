
package main

import (
	"net/http"
	"net/http/httptest"
	"os"
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
					"install": {
						"30d": {"rancher": 1234},
						"90d": {"rancher": 5678},
						"365d": {"rancher": 91011}
					}
				}
			}`))
		} else if strings.HasPrefix(r.URL.Path, "/api/formula") {
			w.Write([]byte(`{
				"analytics": {
					"install": {
						"30d": {"lima": 4321},
						"90d": {"lima": 8765},
						"365d": {"lima": 11019}
					}
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
	if caskAnalytics.Analytics.Install.Downloads30d["rancher"] != 1234 {
		t.Errorf("Expected 30d count for cask to be 1234, but got %d", caskAnalytics.Analytics.Install.Downloads30d["rancher"])
	}

	formulaPkg := homebrewPackage{Name: "lima", Type: "formula"}
	formulaAnalytics, err := fetchAnalyticsData(formulaPkg)
	if err != nil {
		t.Fatalf("fetchAnalyticsData for formula failed: %v", err)
	}
	if formulaAnalytics.Analytics.Install.Downloads30d["lima"] != 4321 {
		t.Errorf("Expected 30d count for formula to be 4321, but got %d", formulaAnalytics.Analytics.Install.Downloads30d["lima"])
	}
}

func totalDownloads(downloads map[string]int) int {
	total := 0
	for _, count := range downloads {
		total += count
	}
	return total
}
