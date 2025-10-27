# Rancher Desktop Download Counter Project Plan

This document outlines the plan to create a system for tracking the download counts of Rancher Desktop release assets over time.

## 1. Goal

The primary goal is to collect daily download statistics for Rancher Desktop release assets to enable trend analysis. Since GitHub only provides cumulative counts, this project will create a historical record of those counts.

## 2. High-Level Approach

A Go program will be executed daily to fetch download counts from the GitHub API. The results will be stored in CSV files. A GitHub Action will be used to automate this process.

## 3. Git Branching Model

The project uses a two-branch model to separate the application code from the data.

-   **`main` branch:** Contains the Go program (`main.go`) and the GitHub Actions workflow.
-   **`data` branch:** Contains the collected download data in CSV files and the daily statistics file (`daily_stats.txt`).

## 4. Implementation Details

### 4.1. Data Fetching and Asset Selection

-   **Language:** The data-fetching script will be written in Go.
-   **Source:** The script will target the `rancher-sandbox/rancher-desktop` repository and fetch all of its releases, not just the latest one.
-   **Asset Filtering:** For each release, the script will only track release assets that have a corresponding `.sha512sum` file in the same release. This focuses the data collection on the main distributable artifacts and filters out source code archives and other assets.

### 4.2. Data Storage

-   **Directory:** All data will be stored in the `data/` directory on the `data` branch.
-   **File Structure:** A separate CSV file will be created for each main asset (e.g., `data/rancher-desktop.msi.csv`).
-   **File Content:** Each CSV file will contain three columns:
    -   `date`: The date of the data capture (YYYY-MM-DD).
    -   `asset_downloads`: The cumulative download count for the main asset file.
    -   `sha512sum_downloads`: The cumulative download count for the corresponding `.sha512sum` file.

This structure allows for easy manual inspection of the data for a specific asset while also capturing the relationship between the asset and its checksum file downloads.

## 5. Development Phases

This project was developed in several phases. The initial development is complete, and the project is now in a maintenance phase.

1.  **Phase 1: Local Script Development**
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
        - Check out the `main` branch.
        - Create a git worktree for the `data` branch in a `data-worktree` directory.
        - Run the Go script from the `data-worktree` directory.
        - Commit any changes in the worktree back to the `data` branch.

## 6. Phase 3: Enhanced Commit Messages with Daily Statistics

To make the daily commit messages more informative, the Go script has been enhanced to generate statistics about the daily download counts. These statistics are used in the body of the commit message.

### 6.1. Statistics to Collect

-   **Total Daily Downloads:** The script calculates the total number of downloads for all main assets (excluding `.sha512sum` files) for the current day.
-   **Top 10 Assets by Daily Downloads:** The script identifies the top 10 assets with the highest download counts for the current day.

### 6.2. Implementation Details

1.  **Go Script Modifications:**
    - The `recordDownloadData` function was updated to calculate the daily download count by subtracting the previous day's total downloads.
    - The `main` function was updated to collect these daily download counts.
    - The `generateStatistics` function was updated to work with daily download counts.
    - The script writes these statistics to a new file, `daily_stats.txt`, in the `data` branch.

2.  **GitHub Actions Workflow Modifications:**
    - The workflow is updated to read the content of `daily_stats.txt` from the worktree.
    - The content of the file is used as the body of the commit message.

## 7. Timezone and Schedule

To ensure that the daily download data is captured for the correct day, the timezone for recording data has been updated.

-   **Schedule:** The GitHub Actions workflow is scheduled to run at 8:00 UTC.
-   **Timezone:** The Go script now uses the `Pacific/Niue` timezone (UTC-11) to record the date for the download counts. This ensures that all downloads for a given day are attributed to the same day, even if the GitHub Action is delayed by several hours.

## 8. Phase 4: Homebrew Analytics

To gain more insight into the distribution of Rancher Desktop, this phase will focus on collecting and storing Homebrew installation analytics.

### 8.1. Goal

The goal of this phase is to track the installation trends of Rancher Desktop on Homebrew over time.

### 8.2. Data Source

The analytics data will be fetched from the Homebrew API. The data is available in a JSON file for each cask. For Rancher Desktop, the URL is `https://formulae.brew.sh/api/cask/rancher.json`.

The API provides the total number of installations over sliding windows of 30, 90, and 365 days. It is not possible to derive daily download numbers from this data. The data is also aggregated and does not provide a breakdown by platform or architecture.

### 8.3. Implementation Details

1.  **Go Program:**
    -   A Go program, `cmd/brew/main.go`, fetches and processes the Homebrew analytics data.
    -   The program fetches the JSON data from the Homebrew API.
    -   It parses the JSON and extracts the 30, 90, and 365-day installation counts.

2.  **Data Storage:**
    -   The data is stored in a new CSV file, `data/homebrew_analytics.csv`.
    -   The CSV file has the following columns:
        -   `date`: The date of data capture (YYYY-MM-DD).
        -   `30d`: The total installations in the last 30 days.
        -   `90d`: The total installations in the last 90 days.
        -   `365d`: The total installations in the last 365 days.

3.  **GitHub Actions Workflow Modifications:**
    -   The existing GitHub Actions workflow runs the `cmd/brew/main.go` program daily.
    -   The workflow commits the updated `data/homebrew_analytics.csv` file to the `data` branch.
