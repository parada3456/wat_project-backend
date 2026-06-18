CREATE TABLE job_overall_ratings (
    rating_summary_id TEXT PRIMARY KEY,
    job_id TEXT NOT NULL UNIQUE REFERENCES job_postings(job_id) ON DELETE CASCADE,
    overall_rate DECIMAL(3,2) NOT NULL DEFAULT 0.0,
    agency_rate DECIMAL(3,2) NOT NULL DEFAULT 0.0,
    job_rate DECIMAL(3,2) NOT NULL DEFAULT 0.0,
    coworkers_rate DECIMAL(3,2) NOT NULL DEFAULT 0.0,
    town_rate DECIMAL(3,2) NOT NULL DEFAULT 0.0,
    hours_rate DECIMAL(3,2) NOT NULL DEFAULT 0.0,
    housing_rate DECIMAL(3,2) NOT NULL DEFAULT 0.0,
    second_job_feasibility_rate DECIMAL(3,2) NOT NULL DEFAULT 0.0,
    overtime_availability_rate DECIMAL(3,2) NOT NULL DEFAULT 0.0,
    review_count INT NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
