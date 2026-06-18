# Execution Plan: Work and Travel Agency Scraper

## Phase 1: Setup and Initialization
1. Initialize browser session via agent.
2. Load configuration and mapping dictionaries from `SPEC.md`.
3. Prepare empty data stores for `job_posting` and `job_housing`.

## Phase 2: Scrape Acadex Agency
1. **Navigate to List Page:** `https://www.acadexthailand.com/program/work-and-travel-summer/`
2. **Pagination/Scrolling:** Scroll or click through pagination to load all available job listings.
3. **URL Extraction:** Extract all individual job detail URLs from the grid/list. *Hint: Also capture the group category (A, X, Y, Z) from the card or URL slug if available here.*
4. **Detail Page Traversal:** For each extracted URL:
    * Navigate to the job detail page.
    * Parse `employer_title`, `location_city`, `location_state`.
    * Parse `position` and `salary` (extract min/max).
    * Parse housing info: `weekly_rate`, `deposit`, `transportation` details.
    * Parse start dates to determine `range_min_start_date` and `range_max_start_date`.
5. **Data Transformation (Acadex):**
    * Apply Group Location mapping (A -> Rank S, etc.).
    * Generate UUIDs for `job_id` and `housing_id`.
    * Append transformed data to the local data stores.

## Phase 3: Scrape iHappy Agency
1. **Navigate to List Page:** `https://www.ihappyeducation.com/job-location-summer/`
2. **Pagination/Scrolling:** Ensure all job cards are loaded.
3. **URL & Badge Extraction:** Extract all individual job detail URLs. Extract the package tier (Signature, Prestige, etc.) from the card badges or headers.
4. **Detail Page Traversal:** For each extracted URL:
    * Navigate to the job detail page.
    * Locate the main data table.
    * Extract `employer_title`, `location_city`, `location_state`.
    * Extract `position` and explicitly grab `position_type` from the table rows.
    * Parse `salary` into min/max numbers.
    * Extract housing `weekly_rate`, `deposit`, and `transportation`.
    * Extract start dates for min/max boundaries.
5. **Data Transformation (iHappy):**
    * Apply Group Location mapping (Signature -> Rank S, etc.).
    * Generate UUIDs for `job_id` and `housing_id`.
    * Append transformed data to the local data stores.

## Phase 4: Output Generation
1. Validate the accumulated data against the types specified in `SPEC.md` (ensure numbers are numeric, dates are formatted, booleans are correct).
2. Export the data into two files:
    * `job_posting.json` (or `.csv`)
    * `job_housing.json` (or `.csv`)
3. Ensure foreign keys (`job_id` in `job_housing`) correctly match the primary keys in `job_posting`.
4. Terminate browser session and return success.