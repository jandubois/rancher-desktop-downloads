
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
		w.Write([]byte(`{
			"analytics": {
				"install": {
					"30d": {
						"rancher": 1234
					},
					"90d": {
						"rancher": 5678
					},
					"365d": {
						"rancher": 91011
					}
				}
			}
		}`))
	}))
	defer server.Close()

	// Override the API URL to point to the mock server
	originalURL := brewAPIURL
	brewAPIURL = server.URL
	defer func() { brewAPIURL = originalURL }()

	// Run the tests
	os.Exit(m.Run())
}

func TestFetchCaskData(t *testing.T) {
	cask, err := fetchCaskData()
	if err != nil {
		t.Fatalf("fetchCaskData failed: %v", err)
	}

	if cask.Analytics.Install.ThirtyDays["rancher"] != 1234 {
		t.Errorf("Expected 30d count to be 1234, but got %d", cask.Analytics.Install.ThirtyDays["rancher"])
	}
	if cask.Analytics.Install.NinetyDays["rancher"] != 5678 {
		t.Errorf("Expected 90d count to be 5678, but got %d", cask.Analytics.Install.NinetyDays["rancher"])
	}
	if cask.Analytics.Install.ThreeHundredSixtyFiveDays["rancher"] != 91011 {
		t.Errorf("Expected 365d count to be 91011, but got %d", cask.Analytics.Install.ThreeHundredSixtyFiveDays["rancher"])
	}
}
