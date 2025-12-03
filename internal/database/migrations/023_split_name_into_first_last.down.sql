-- Rollback migration: combine first_name and last_name back into name column
-- Migrate data back: combine first_name and last_name into name column
UPDATE users
SET name = TRIM(CONCAT_WS(' ', first_name, last_name))
WHERE name IS NULL OR name = '';

-- Drop indexes
DROP INDEX IF EXISTS idx_users_first_name;
DROP INDEX IF EXISTS idx_users_last_name;

-- Drop the new columns
ALTER TABLE users
DROP COLUMN IF EXISTS first_name,
DROP COLUMN IF EXISTS last_name;
