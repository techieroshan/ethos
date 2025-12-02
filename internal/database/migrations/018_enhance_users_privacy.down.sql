-- Remove privacy control columns from users table
ALTER TABLE users
DROP COLUMN IF EXISTS opt_outs,
DROP COLUMN IF EXISTS anonymized_at,
DROP COLUMN IF EXISTS delete_requested_at;
