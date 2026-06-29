ALTER TABLE users ADD COLUMN first_name TEXT NOT NULL DEFAULT '';
ALTER TABLE users ADD COLUMN last_name TEXT NOT NULL DEFAULT '';

UPDATE users u
SET first_name = p.first_name,
    last_name = p.last_name
FROM profiles p
WHERE u.user_id = p.user_id;

ALTER TABLE users ALTER COLUMN first_name DROP DEFAULT;
ALTER TABLE users ALTER COLUMN last_name DROP DEFAULT;

ALTER TABLE profiles DROP COLUMN first_name;
ALTER TABLE profiles DROP COLUMN last_name;
