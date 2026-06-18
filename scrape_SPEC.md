# Scraping Specification for Work and Travel Agencies

## 1. Objective
Extract Job Posting and Job Housing data from Acadex and iHappy agency websites and format it to match the target SQL schema. using golang

## 2. Target Schema (JSON Representation of SQL)

### 2.1 Job Posting (`job_posting`)
* **job_id** (String): UUID or generated hash based on `source_url`.
* **agency_name** (String): "Acadex" or "iHappy".
* **employer_title** (String): Name of the employer/company.
* **position** (String): Job role/title.
* **position_type** (String): E.g., Resort worker, Fast Food, Retail. 
* **location_city** (String): Extracted from location details.
* **location_state** (String): Extracted from location details.
* **group_location** (String): **MUST** be mapped to internal ranks (See Section 3).
* **us_sponsor** (Boolean): Default `true`.
* **salary_range_min** (Number): Parsed from hourly wage text.
* **salary_range_max** (Number): Parsed from hourly wage text (same as min if flat rate).
* **available_slots** (Integer): Number of positions open.
* **description** (String): General job description or requirements.
* **source_url** (String): URL of the job detail page.
* **scrape_at** (Timestamp): ISO 8601 UTC timestamp of scraping.

### 2.2 Job Housing (`job_housing`)
* **housing_id** (String): UUID or generated hash.
* **job_id** (String): Foreign key to `job_posting`.
* **description** (String): Details about the housing arrangement.
* **weekly_rate** (Number): Parsed rent cost per week.
* **deposit** (Number): Parsed housing deposit amount.
* **transportation** (String): Commute details (e.g., walk, bike, bus).
* **range_min_start_date** (Date): Earliest allowed start date (YYYY-MM-DD).
* **range_max_start_date** (Date): Latest allowed start date (YYYY-MM-DD).

---

## 3. Business Logic & Mapping Rules

### 3.1 Group Location Mapping
You must strictly map the agency-specific groupings to the internal ranking system:

**Acadex Mapping:**
* Group A -> `Rank S`
* Group X -> `Rank A`
* Group Y -> `Rank B`
* Group Z -> `Rank C`
*(Note: Look for these group indicators in the listing tags, filters, or the URL slug itself).*

**iHappy Mapping:**
* Signature -> `Rank S`
* Prestige -> `Rank A`
* Super Premium -> `Rank B`
* Premium -> `Rank C`
* Promotion -> `Rank D`
*(Note: Look for these identifiers in the badges, headers, or categorical tags on the list/detail pages).*

### 3.2 Position Type Extraction
* **Acadex:** Extract from the job category filters on the list page or breadcrumbs/tags on the detail page.
* **iHappy:** Extract directly from the structured information table found on the individual job detail page.

### 3.3 Data Cleaning Rules
* **Currency:** Strip `$` signs and convert to numeric values for salary, deposit, and weekly rates.
* **Dates:** Convert text strings like "May 15th - June 1st" into discrete `range_min_start_date` and `range_max_start_date` in `YYYY-MM-DD` format. Assume the target year is the one mentioned in the title/URL (e.g., 2027).