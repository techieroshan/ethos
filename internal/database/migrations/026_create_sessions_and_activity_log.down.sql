-- Down migration for sessions and activity log
-- +migrate Down

DROP INDEX IF EXISTS idx_user_context_switches_timestamp;
DROP INDEX IF EXISTS idx_user_context_switches_user_id;
DROP INDEX IF EXISTS idx_organization_activity_log_action;
DROP INDEX IF EXISTS idx_organization_activity_log_created_at;
DROP INDEX IF EXISTS idx_organization_activity_log_user_id;
DROP INDEX IF EXISTS idx_organization_activity_log_org_id;
DROP INDEX IF EXISTS idx_user_sessions_revoked_at;
DROP INDEX IF EXISTS idx_user_sessions_expires_at;
DROP INDEX IF EXISTS idx_user_sessions_organization_id;
DROP INDEX IF EXISTS idx_user_sessions_user_id;

DROP TABLE IF EXISTS user_context_switches;
DROP TABLE IF EXISTS organization_activity_log;
DROP TABLE IF EXISTS user_sessions;
