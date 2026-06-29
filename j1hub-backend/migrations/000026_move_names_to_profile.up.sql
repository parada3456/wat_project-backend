ALTER TABLE profiles ADD COLUMN first_name TEXT NOT NULL DEFAULT '';
ALTER TABLE profiles ADD COLUMN last_name TEXT NOT NULL DEFAULT '';

UPDATE profiles p
SET first_name = u.first_name,
    last_name = u.last_name
FROM users u
WHERE p.user_id = u.user_id;

ALTER TABLE profiles ALTER COLUMN first_name DROP DEFAULT;
ALTER TABLE profiles ALTER COLUMN last_name DROP DEFAULT;

ALTER TABLE users DROP COLUMN first_name;
ALTER TABLE users DROP COLUMN last_name;
