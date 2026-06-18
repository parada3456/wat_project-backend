CREATE TABLE user_missions (
    user_mission_id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    mission_id TEXT NOT NULL REFERENCES missions(mission_id) ON DELETE CASCADE,
    status TEXT NOT NULL CHECK (status IN ('Not_Started', 'In_Progress', 'Pending_Verification', 'Completed', 'Overdue')),
    calculated_due_date TIMESTAMPTZ NOT NULL,
    proof_url TEXT,
    proof_submitted_at TIMESTAMPTZ,
    verified_at TIMESTAMPTZ,
    verified_by TEXT,
    base_points_earned INT NOT NULL DEFAULT 0,
    speed_bonus_points INT NOT NULL DEFAULT 0,
    streak_bonus_points INT NOT NULL DEFAULT 0,
    first_completer_bonus_points INT NOT NULL DEFAULT 0,
    total_points_earned INT NOT NULL DEFAULT 0,
    rewarded_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
