# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a download statistics tracking system for Rancher Desktop and related desktop virtualization tools. The system collects daily download counts from three sources (GitHub releases, Homebrew analytics, and Chocolatey packages), storing historical data in CSV files on a separate `data` branch.

## Architecture

### Branch Model
- **`main` branch**: Contains Go programs and GitHub Actions workflow
- **`data` branch**: Contains collected CSV files and daily statistics (updated daily by automation)

The workflow uses `git worktree` to check out the data branch into `data-worktree/` for updates.

### Data Collection Programs

Three independent Go programs collect statistics from different sources:

1. **`cmd/github/main.go`**: Fetches GitHub release asset download counts
   - Tracks assets paired with their `.sha512sum` checksum files
   - Filters out source archives and unpaired assets
   - Repository: `rancher-sandbox/rancher-desktop`
   - Uses GitHub API (requires `GITHUB_TOKEN` or `GH_TOKEN` env var)

2. **`cmd/brew/main.go`**: Fetches Homebrew installation analytics
   - Tracks both casks and formulae (rancher, docker-desktop, lima, colima, virtualbox, vagrant, multipass, finch, utm, qemu, podman-desktop)
   - Records 30d/90d/365d rolling installation counts
   - Uses Homebrew's public API at `formulae.brew.sh`
   - No authentication required

3. **`cmd/chocolatey/main.go`**: Fetches Chocolatey package download counts
   - Scrapes package pages at `community.chocolatey.org`
   - Tracks rancher-desktop, docker-desktop, podman-desktop
   - Records total downloads and version-specific downloads
   - Stores data in `data/chocolatey/` subdirectory
   - No authentication required

### Shared Data Package

The `internal/data` package provides:
- `RecordDownloadData()`: Writes GitHub download data to CSV files with cumulative counts. Returns daily download delta for statistics generation.
- `RecordBrewAnalytics()`: Writes Homebrew analytics to CSV files with rolling window counts
- `RecordChocolateyTotal()`: Writes total Chocolatey download counts to CSV files
- `RecordChocolateyVersion()`: Writes version-specific Chocolatey download counts to CSV files (creates new file per version)
- `GenerateStatistics()`: Creates `daily_stats.txt` with daily download totals and top 10 assets
- Timezone handling: Uses `Pacific/Niue` (UTC-11) to ensure consistent date recording even with GitHub Actions delays

### CSV File Format

**GitHub files**: `data/<asset-name>.csv`
```csv
date,asset_downloads,sha512sum_downloads
2024-01-01,1000,950
```
- Cumulative download counts for each asset and its `.sha512sum` file

**Homebrew files**: `data/<package-name>-<type>.csv`
```csv
date,30d,90d,365d
2024-01-01,5000,12000,50000
```
- Rolling window installation counts (30/90/365 days)

**Chocolatey files**: `data/chocolatey/<package>-total.csv` and `data/chocolatey/<package>-<version>.csv`
```csv
date,total_downloads
2024-01-01,212699
```
```csv
date,version_downloads
2024-01-01,291
```
- Total downloads tracked in `*-total.csv` files
- Version-specific downloads in separate files (new file created when version changes)

## Development Commands

### Run tests
```bash
go test -v ./...
```

### Run individual collector locally
```bash
# GitHub (requires GITHUB_TOKEN or GH_TOKEN env var)
go run cmd/github/main.go

# Homebrew (no auth required)
go run cmd/brew/main.go

# Chocolatey (no auth required)
go run cmd/chocolatey/main.go
```

### Run a single test
```bash
# Test a specific function
go test -v -run TestFetchReleases ./cmd/github

# Test a specific package
go test -v ./internal/data
```

### Trigger workflow manually
The GitHub Actions workflow can be triggered via `workflow_dispatch` in addition to the daily 8:00 UTC schedule.

## Key Implementation Details

- **Daily vs cumulative tracking**: GitHub and Chocolatey provide cumulative counts, so `RecordDownloadData()` and `RecordChocolateyTotal()` calculate daily deltas by comparing to previous day's value. Homebrew provides rolling windows directly.
- **Same-day overwrites**: If run multiple times in the same day (Niue timezone), the last record for that day is replaced rather than appended.
- **Statistics generation**: The github program generates `daily_stats.txt` which becomes the commit message body in the automated workflow.
- **Asset pairing logic**: GitHub assets are only tracked if they have a matching `.sha512sum` file (see `pairAssets()` in cmd/github/main.go:108).
- **Chocolatey web scraping**: Uses `goquery` library to parse HTML from Chocolatey package pages. Extracts data from `#totalDownloadCount`, `#currentDownloadCount` elements, and version from Nuspec metadata (see cmd/chocolatey/main.go:65).
- **Version-specific tracking**: Chocolatey creates a new CSV file whenever the latest version changes, allowing historical tracking of per-version adoption.
- **Testing approach**: Tests use `httptest` to mock API/web servers and `t.TempDir()` for isolated file operations. The `data.DataDir` variable is overridden in tests for proper isolation.
