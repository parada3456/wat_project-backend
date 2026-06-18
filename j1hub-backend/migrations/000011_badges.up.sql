CREATE TABLE badges (
    badge_id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    trigger_type TEXT NOT NULL CHECK (trigger_type IN ('Speed', 'Streak', 'First_Completer', 'Phase_Complete', 'Manual')),
    icon_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
