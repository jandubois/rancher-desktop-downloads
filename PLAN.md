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
