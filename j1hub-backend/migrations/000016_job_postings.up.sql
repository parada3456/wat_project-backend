CREATE TABLE job_postings (
    job_id TEXT PRIMARY KEY,
    agency_name TEXT NOT NULL,
    employer_title TEXT NOT NULL,
    position TEXT NOT NULL,
    position_type TEXT NOT NULL,
    location_city TEXT NOT NULL,
    location_state TEXT NOT NULL,
    group_location TEXT,
    us_sponsor BOOLEAN NOT NULL DEFAULT true,
    salary_range_min DECIMAL(12,2),
    salary_range_max DECIMAL(12,2),
    available_slots INT NOT NULL DEFAULT 0,
    description TEXT,
    source_url TEXT,
    scrape_at TIMESTAMPTZ,
    posted_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
