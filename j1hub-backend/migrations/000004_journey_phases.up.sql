CREATE TABLE journey_phases (
    phase_id TEXT PRIMARY KEY,
    phase_number INT UNIQUE NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE users ADD CONSTRAINT fk_current_phase FOREIGN KEY (current_phase_id) REFERENCES journey_phases(phase_id);
