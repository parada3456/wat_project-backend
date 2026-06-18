CREATE TABLE user_carts (
    cart_id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    job_id TEXT NOT NULL REFERENCES job_postings(job_id) ON DELETE CASCADE,
    status TEXT NOT NULL CHECK (status IN ('Saved', 'Viewed', 'Applied', 'Removed')),
    added_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
