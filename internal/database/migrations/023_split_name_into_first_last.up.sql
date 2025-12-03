-- Add first_name and last_name columns to users table
ALTER TABLE users
ADD COLUMN IF NOT EXISTS first_name VARCHAR(100),
ADD COLUMN IF NOT EXISTS last_name VARCHAR(100);

-- Migrate existing data: split name column into first_name and last_name
-- If name has a space, split it; otherwise put it in first_name
UPDATE users
SET 
    first_name = CASE 
        WHEN name IS NULL THEN ''
        WHEN position(' ' in name) > 0 THEN substring(name from 1 for position(' ' in name) - 1)
        ELSE name
    END,
    last_name = CASE 
        WHEN name IS NULL THEN ''
        WHEN position(' ' in name) > 0 THEN substring(name from position(' ' in name) + 1)
        ELSE ''
    END
WHERE first_name IS NULL AND last_name IS NULL;

-- Make the new columns NOT NULL with DEFAULT empty string for backward compatibility
ALTER TABLE users
ALTER COLUMN first_name SET NOT NULL,
ALTER COLUMN first_name SET DEFAULT '',
ALTER COLUMN last_name SET NOT NULL,
ALTER COLUMN last_name SET DEFAULT '';

-- Create indexes for new columns if needed
CREATE INDEX IF NOT EXISTS idx_users_first_name ON users(first_name);
CREATE INDEX IF NOT EXISTS idx_users_last_name ON users(last_name);
