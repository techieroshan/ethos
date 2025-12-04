-- Down migration for users table enhancements
-- +migrate Down

ALTER TABLE users
DROP CONSTRAINT IF EXISTS check_locked_until_after_now;

DROP INDEX IF EXISTS idx_users_current_organization_id;
DROP INDEX IF EXISTS idx_users_last_login_at;
DROP INDEX IF EXISTS idx_users_two_factor_enabled;
DROP INDEX IF EXISTS idx_users_account_status;

ALTER TABLE users
DROP COLUMN IF EXISTS preferences,
DROP COLUMN IF EXISTS current_organization_id,
DROP COLUMN IF EXISTS locked_until,
DROP COLUMN IF EXISTS login_attempts,
DROP COLUMN IF EXISTS two_factor_secret,
DROP COLUMN IF EXISTS two_factor_enabled,
DROP COLUMN IF EXISTS last_login_at,
DROP COLUMN IF EXISTS account_status,
DROP COLUMN IF EXISTS timezone,
DROP COLUMN IF EXISTS phone_number,
DROP COLUMN IF EXISTS avatar_url;
