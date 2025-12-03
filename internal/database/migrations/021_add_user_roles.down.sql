-- Remove user roles system
-- +migrate Down

-- Drop triggers
DROP TRIGGER IF EXISTS update_user_roles_updated_at ON user_roles;

-- Drop tables in reverse order
DROP TABLE IF EXISTS user_role_assignments;
DROP TABLE IF EXISTS user_roles;
