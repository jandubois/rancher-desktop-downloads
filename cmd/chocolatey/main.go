package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"rancher-desktop-downloads/internal/data"
)

var (
	packages = []string{
		"rancher-desktop",
		"docker-desktop",
		"podman-desktop",
	}
	baseURL = "https://community.chocolatey.org/packages/"
)

// PackageStats represents the statistics for a Chocolatey package.
type PackageStats struct {
	TotalDownloads   int
	Version          string
	VersionDownloads int
}

func main() {
	chocolateyDataDir := "data/chocolatey"
	if err := os.MkdirAll(chocolateyDataDir, 0755); err != nil {
		fmt.Printf("Error creating chocolatey data directory: %v\n", err)
		os.Exit(1)
	}

	for _, pkg := range packages {
		fmt.Printf("Fetching Chocolatey stats for %s...\n", pkg)
		stats, err := fetchPackageStats(pkg)
		if err != nil {
			fmt.Printf("Error fetching stats for %s: %v\n", pkg, err)
			continue
		}

		fmt.Printf("  Total downloads: %d\n", stats.TotalDownloads)
		fmt.Printf("  Latest version: %s (%d downloads)\n", stats.Version, stats.VersionDownloads)

		// Record total downloads
		if err := data.RecordChocolateyTotal(pkg, stats.TotalDownloads); err != nil {
			fmt.Printf("Error recording total downloads for %s: %v\n", pkg, err)
			continue
		}

		// Record version-specific downloads
		if err := data.RecordChocolateyVersion(pkg, stats.Version, stats.VersionDownloads); err != nil {
			fmt.Printf("Error recording version downloads for %s: %v\n", pkg, err)
			continue
		}
	}

	fmt.Println("Chocolatey stats collection complete.")
}

// fetchPackageStats fetches statistics for a Chocolatey package by scraping the website.
func fetchPackageStats(packageName string) (*PackageStats, error) {
	url := baseURL + packageName
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	stats := &PackageStats{}

	// Extract total downloads
	totalDownloadsText := doc.Find("#totalDownloadCount").Text()
	totalDownloads, err := parseNumber(totalDownloadsText)
	if err != nil {
		return nil, fmt.Errorf("failed to parse total downloads: %w", err)
	}
	stats.TotalDownloads = totalDownloads

	// Extract current version downloads
	currentDownloadsText := doc.Find("#currentDownloadCount").Text()
	currentDownloads, err := parseNumber(currentDownloadsText)
	if err != nil {
		return nil, fmt.Errorf("failed to parse current version downloads: %w", err)
	}
	stats.VersionDownloads = currentDownloads

	// Extract version number from the Nuspec reference
	versionFound := false
	doc.Find("p[title]").Each(func(i int, s *goquery.Selection) {
		title, exists := s.Attr("title")
		if exists && strings.Contains(title, "Nuspec reference: <version>") {
			version := extractVersionFromNuspec(title)
			if version != "" {
				stats.Version = version
				versionFound = true
			}
		}
	})

	if !versionFound {
		return nil, fmt.Errorf("failed to find version number")
	}

	return stats, nil
}

// parseNumber converts a string like "212,699" to an integer 212699.
func parseNumber(s string) (int, error) {
	// Remove commas and whitespace
	cleaned := strings.ReplaceAll(strings.TrimSpace(s), ",", "")
	return strconv.Atoi(cleaned)
}

// extractVersionFromNuspec extracts the version from a string like "Nuspec reference: <version>1.20.1</version>".
func extractVersionFromNuspec(nuspecTitle string) string {
	re := regexp.MustCompile(`<version>([^<]+)</version>`)
	matches := re.FindStringSubmatch(nuspecTitle)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}
