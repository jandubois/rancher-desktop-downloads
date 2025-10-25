# Rancher Desktop Download Counter Project Plan

This document outlines the plan to create a system for tracking the download counts of Rancher Desktop release assets over time.

## 1. Goal

The primary goal is to collect daily download statistics for Rancher Desktop release assets to enable trend analysis. Since GitHub only provides cumulative counts, this project will create a historical record of those counts.

## 2. High-Level Approach

A Go program will be executed daily to fetch download counts from the GitHub API. The results will be stored in CSV files within this repository. A GitHub Action will be used to automate this process.

## 3. Implementation Details

### 3.1. Data Fetching and Asset Selection

- **Language:** The data-fetching script will be written in Go.
- **Source:** The script will target the `rancher-sandbox/rancher-desktop` repository and fetch all of its releases, not just the latest one.
- **Asset Filtering:** For each release, the script will only track release assets that have a corresponding `.sha512sum` file in the same release. This focuses the data collection on the main distributable artifacts and filters out source code archives and other assets.

### 3.2. Data Storage

- **Directory:** All data will be stored in the `data/` directory.
- **File Structure:** A separate CSV file will be created for each main asset (e.g., `data/rancher-desktop.msi.csv`).
- **File Content:** Each CSV file will contain three columns:
    - `date`: The date of the data capture (YYYY-MM-DD).
    - `asset_downloads`: The cumulative download count for the main asset file.
    - `sha512sum_downloads`: The cumulative download count for the corresponding `.sha512sum` file.

This structure allows for easy manual inspection of the data for a specific asset while also capturing the relationship between the asset and its checksum file downloads.

## 4. Development Phases

1.  **Phase 1: Local Script Development (Current Phase)**
    - Create this `PLAN.md` document.
    - Develop the Go script (`main.go`) to implement the data fetching and file writing logic.
    - Write unit tests for the file writing logic to ensure correctness.
    - Test the script locally to ensure it correctly fetches data and creates the CSV files.
    - Create a `.gitignore` file.

2.  **Phase 2: GitHub Action Automation**
    - Create a GitHub Actions workflow file in `.github/workflows/`.
    - The workflow will trigger on a daily schedule.
    - The workflow will:
        - Set up the Go environment.
        - Run the Go script.
        - Commit any changes to the `data/` directory back to the repository.

## 5. Phase 3: Enhanced Commit Messages with Daily Statistics

To make the daily commit messages more informative, the Go script has been enhanced to generate statistics about the daily download counts. These statistics are used in the body of the commit message.

### 5.1. Statistics to Collect

- **Total Daily Downloads:** The script calculates the total number of downloads for all main assets (excluding `.sha512sum` files) for the current day.
- **Top 10 Assets by Daily Downloads:** The script identifies the top 10 assets with the highest download counts for the current day.

### 5.2. Implementation Details

1.  **Go Script Modifications:**
    - The `recordDownloadData` function was updated to calculate the daily download count by subtracting the previous day's total downloads.
    - The `main` function was updated to collect these daily download counts.
    - The `generateStatistics` function was updated to work with daily download counts.
    - The script writes these statistics to a new file, `daily_stats.txt`.

2.  **GitHub Actions Workflow Modifications:**
    - The workflow is updated to read the content of `daily_stats.txt`.
    - The content of the file is used as the body of the commit message.

## 6. Timezone and Schedule

To better capture download data from the Americas, the daily execution time and the timezone used for recording data have been updated.

- **Schedule:** The GitHub Actions workflow is scheduled to run at midnight Pacific Time (`0 8 * * *` UTC).
- **Timezone:** The Go script now uses the `America/Los_Angeles` timezone to record the date for the download counts. This ensures that all downloads for a given day in the Americas are attributed to the same day.
