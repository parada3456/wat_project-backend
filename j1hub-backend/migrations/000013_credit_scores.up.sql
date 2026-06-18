CREATE TABLE credit_scores (
    credit_id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL UNIQUE REFERENCES users(user_id) ON DELETE CASCADE,
    current_score INT NOT NULL DEFAULT 100,
    last_updated TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
