CREATE TABLE friendships (
    friendship_id TEXT PRIMARY KEY,
    user_id_1 TEXT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    user_id_2 TEXT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    status TEXT NOT NULL CHECK (status IN ('Pending', 'Accepted', 'Blocked')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id_1, user_id_2)
);
-- App enforces user_id_1 < user_id_2 for canonical order
