package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Mock the Chocolatey server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
<!DOCTYPE html>
<html>
<head>
    <title>Chocolatey Software | Test Package 1.0.0</title>
</head>
<body>
    <p title="Nuspec reference: <version>1.0.0</version>" class="mb-0">1.0.0 | <strong>Updated: </strong>01 Nov 2025</p>
    <h3 id="totalDownloadCount" class="mt-0 mb-3">10,000</h3>
    <p class="mb-0"><strong>Downloads of v 1.0.0:</strong></p>
    <h3 id="currentDownloadCount" class="mt-0 mb-3">500</h3>
</body>
</html>
`))
	}))
	defer server.Close()

	// Override the base URL to point to the mock server
	originalURL := baseURL
	baseURL = server.URL + "/"
	defer func() { baseURL = originalURL }()

	// Run the tests
	os.Exit(m.Run())
}

func TestFetchPackageStats(t *testing.T) {
	stats, err := fetchPackageStats("test-package")
	if err != nil {
		t.Fatalf("fetchPackageStats failed: %v", err)
	}

	if stats.TotalDownloads != 10000 {
		t.Errorf("Expected total downloads to be 10000, but got %d", stats.TotalDownloads)
	}

	if stats.Version != "1.0.0" {
		t.Errorf("Expected version to be 1.0.0, but got %s", stats.Version)
	}

	if stats.VersionDownloads != 500 {
		t.Errorf("Expected version downloads to be 500, but got %d", stats.VersionDownloads)
	}
}

func TestParseNumber(t *testing.T) {
	tests := []struct {
		input    string
		expected int
		hasError bool
	}{
		{"212,699", 212699, false},
		{"10,000", 10000, false},
		{"500", 500, false},
		{"1,234,567", 1234567, false},
		{"  100  ", 100, false},
		{"invalid", 0, true},
	}

	for _, tt := range tests {
		result, err := parseNumber(tt.input)
		if tt.hasError {
			if err == nil {
				t.Errorf("parseNumber(%q) expected error, but got nil", tt.input)
			}
		} else {
			if err != nil {
				t.Errorf("parseNumber(%q) unexpected error: %v", tt.input, err)
			}
			if result != tt.expected {
				t.Errorf("parseNumber(%q) = %d, expected %d", tt.input, result, tt.expected)
			}
		}
	}
}

func TestExtractVersionFromNuspec(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Nuspec reference: <version>1.20.1</version>", "1.20.1"},
		{"Nuspec reference: <version>2.0.0</version>", "2.0.0"},
		{"Nuspec reference: <version>1.0.0-beta</version>", "1.0.0-beta"},
		{"No version here", ""},
	}

	for _, tt := range tests {
		result := extractVersionFromNuspec(tt.input)
		if result != tt.expected {
			t.Errorf("extractVersionFromNuspec(%q) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}
