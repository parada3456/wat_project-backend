CREATE TABLE missions (
    mission_id TEXT PRIMARY KEY,
    phase_id TEXT NOT NULL REFERENCES journey_phases(phase_id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT,
    location TEXT,
    base_points INT NOT NULL DEFAULT 0,
    is_mandatory BOOLEAN NOT NULL DEFAULT true,
    verification_type TEXT NOT NULL CHECK (verification_type IN ('None', 'Upload', 'Admin')),
    due_date_type TEXT NOT NULL CHECK (due_date_type IN ('Relative', 'Fixed')),
    fixed_due_date TIMESTAMPTZ,
    relative_trigger_event TEXT CHECK (relative_trigger_event IN ('arrival_date', 'job_start_date')),
    relative_days_offset INT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
