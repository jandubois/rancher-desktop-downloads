# Rancher Desktop Download Counter Project

This document describes the system for tracking the download counts of Rancher Desktop release assets over time.

## 1. Goal

The primary goal is to collect daily download statistics for Rancher Desktop release assets to enable trend analysis. Since GitHub only provides cumulative counts, this project creates a historical record of those counts.

## 2. System Architecture

### 2.1. Overview

The system consists of two Go programs that are executed daily by a GitHub Actions workflow:

-   `cmd/github/main.go`: Fetches download counts for Rancher Desktop release assets from the GitHub API.
-   `cmd/brew/main.go`: Fetches installation analytics for the Rancher Desktop Homebrew cask from the Homebrew API.

The collected data is stored in CSV files on a separate `data` branch of the repository.

### 2.2. Data Sources

-   **GitHub API**: The system fetches release and asset information from the `rancher-sandbox/rancher-desktop` repository using the GitHub API.
-   **Homebrew API**: The system fetches Homebrew installation analytics from the `https://formulae.brew.sh/api/cask/rancher.json` endpoint.

### 2.3. Data Storage

All data is stored in CSV files in the `data/` directory on the `data` branch.

### 2.4. Automation

A GitHub Actions workflow in `.github/workflows/daily-download-count.yml` automates the process of fetching and storing the data. The workflow is scheduled to run daily.

## 3. Data Schema

### 3.1. GitHub Download Counts

A separate CSV file is created for each release asset (e.g., `data/rancher-desktop.msi.csv`). Each file contains the following columns:

-   `date`: The date of the data capture (YYYY-MM-DD).
-   `asset_downloads`: The cumulative download count for the main asset file.
-   `sha512sum_downloads`: The cumulative download count for the corresponding `.sha512sum` file.

### 3.2. Homebrew Analytics

The Homebrew analytics data is stored in `data/homebrew_analytics.csv`. This file contains the following columns:

-   `date`: The date of data capture (YYYY-MM-DD).
-   `30d`: The total installations in the last 30 days.
-   `90d`: The total installations in the last 90 days.
-   `365d`: The total installations in the last 365 days.

## 4. Implementation Details

### 4.1. GitHub Analytics

The `cmd/github/main.go` program fetches all releases for the `rancher-sandbox/rancher-desktop` repository. For each release, it only tracks release assets that have a corresponding `.sha512sum` file in the same release. This filters out source code archives and other non-essential assets.

### 4.2. Homebrew Analytics

The `cmd/brew/main.go` program fetches the JSON data from the Homebrew API, parses it, and extracts the 30, 90, and 365-day installation counts.

### 4.3. Commit Message Generation

To make the daily commit messages more informative, the `cmd/github/main.go` program generates statistics about the daily download counts. These statistics are written to `daily_stats.txt` and used as the body of the commit message. The statistics include:

-   **Total Daily Downloads**: The total number of downloads for all main assets for the current day.
-   **Top 10 Assets by Daily Downloads**: The top 10 assets with the highest download counts for the current day.

## 5. Operational Details

### 5.1. Schedule

The GitHub Actions workflow is scheduled to run daily at 8:00 UTC.

### 5.2. Timezone

The Go scripts use the `Pacific/Niue` timezone (UTC-11) to record the date for the download counts. This ensures that all downloads for a given day are attributed to the same day, even if the GitHub Action is delayed by several hours.

### 5.3. Git Branching Model

The project uses a two-branch model to separate the application code from the data:

-   **`main` branch:** Contains the Go programs and the GitHub Actions workflow.
-   **`data` branch:** Contains the collected download data in CSV files and the daily statistics file (`daily_stats.txt`).