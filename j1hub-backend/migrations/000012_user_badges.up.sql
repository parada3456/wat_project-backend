CREATE TABLE user_badges (
    user_badge_id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    badge_id TEXT NOT NULL REFERENCES badges(badge_id) ON DELETE CASCADE,
    source_id TEXT,
    earned_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
