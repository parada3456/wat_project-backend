CREATE TABLE job_reviews (
    review_id TEXT PRIMARY KEY,
    job_id TEXT NOT NULL REFERENCES job_postings(job_id) ON DELETE CASCADE,
    user_id TEXT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    rating_stars DECIMAL(3,2) NOT NULL,
    review_text TEXT,
    tips_for_next_generation TEXT,
    score_agency DECIMAL(3,2) NOT NULL,
    score_job DECIMAL(3,2) NOT NULL,
    score_coworkers DECIMAL(3,2) NOT NULL,
    score_town DECIMAL(3,2) NOT NULL,
    score_hours DECIMAL(3,2) NOT NULL,
    score_housing DECIMAL(3,2) NOT NULL,
    score_second_job_feasibility DECIMAL(3,2) NOT NULL,
    score_overtime_availability DECIMAL(3,2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
