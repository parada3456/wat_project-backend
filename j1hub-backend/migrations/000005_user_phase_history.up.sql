CREATE TABLE user_phase_history (
    history_id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    phase_id TEXT NOT NULL REFERENCES journey_phases(phase_id) ON DELETE CASCADE,
    phase_points_earned INT NOT NULL DEFAULT 0,
    entered_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ
);
