CREATE TABLE job_housings (
    housing_id TEXT PRIMARY KEY,
    job_id TEXT NOT NULL REFERENCES job_postings(job_id) ON DELETE CASCADE,
    description TEXT,
    weekly_rate DECIMAL(12,2),
    deposit DECIMAL(12,2),
    transportation TEXT,
    range_min_start_date DATE,
    range_max_start_date DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
