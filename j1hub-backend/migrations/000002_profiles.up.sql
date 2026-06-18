CREATE TABLE profiles (
    profile_id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL UNIQUE REFERENCES users(user_id) ON DELETE CASCADE,
    phone_number TEXT,
    bio TEXT,
    avatar_url TEXT,
    radar_visibility TEXT NOT NULL CHECK (radar_visibility IN ('Show_Anonymous', 'Show_Friends', 'Hidden')),
    current_coordinates GEOMETRY(POINT, 4326),
    location_updated_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ON profiles USING GIST (current_coordinates);
