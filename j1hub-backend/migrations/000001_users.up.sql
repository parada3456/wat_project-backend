CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE users (
    user_id TEXT PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    current_phase_id TEXT, -- FK added later to avoid circular dep
    total_lifetime_points INT NOT NULL DEFAULT 0,
    current_phase_points INT NOT NULL DEFAULT 0,
    mission_streak INT NOT NULL DEFAULT 0,
    arrival_date TIMESTAMPTZ,
    job_start_date TIMESTAMPTZ,
    fcm_token TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
