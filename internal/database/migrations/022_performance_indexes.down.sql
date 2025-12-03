-- Rollback performance optimization indexes and constraints

-- Drop all the indexes created in the up migration
DROP INDEX CONCURRENTLY IF EXISTS idx_users_email;
DROP INDEX CONCURRENTLY IF EXISTS idx_users_created_at;
DROP INDEX CONCURRENTLY IF EXISTS idx_users_email_verified;

DROP INDEX CONCURRENTLY IF EXISTS idx_feedback_author_id;
DROP INDEX CONCURRENTLY IF EXISTS idx_feedback_recipient_id;
DROP INDEX CONCURRENTLY IF EXISTS idx_feedback_created_at;
DROP INDEX CONCURRENTLY IF EXISTS idx_feedback_visibility;
DROP INDEX CONCURRENTLY IF EXISTS idx_feedback_rating;
DROP INDEX CONCURRENTLY IF EXISTS idx_feedback_tenant_author;
DROP INDEX CONCURRENTLY IF EXISTS idx_feedback_tenant_created;
DROP INDEX CONCURRENTLY IF EXISTS idx_feedback_content_fts;

DROP INDEX CONCURRENTLY IF EXISTS idx_comments_feedback_id;
DROP INDEX CONCURRENTLY IF EXISTS idx_comments_author_id;
DROP INDEX CONCURRENTLY IF EXISTS idx_comments_created_at;
DROP INDEX CONCURRENTLY IF EXISTS idx_comments_tenant_feedback;

DROP INDEX CONCURRENTLY IF EXISTS idx_reactions_feedback_id;
DROP INDEX CONCURRENTLY IF EXISTS idx_reactions_user_id;
DROP INDEX CONCURRENTLY IF EXISTS idx_reactions_type;
DROP INDEX CONCURRENTLY IF EXISTS idx_reactions_tenant_user;

DROP INDEX CONCURRENTLY IF EXISTS idx_user_role_assignments_user_id;
DROP INDEX CONCURRENTLY IF EXISTS idx_user_role_assignments_role_id;
DROP INDEX CONCURRENTLY IF EXISTS idx_user_role_assignments_active;
DROP INDEX CONCURRENTLY IF EXISTS idx_user_role_assignments_user_active;

DROP INDEX CONCURRENTLY IF EXISTS idx_user_tenant_memberships_user_id;
DROP INDEX CONCURRENTLY IF EXISTS idx_user_tenant_memberships_tenant_id;
DROP INDEX CONCURRENTLY IF EXISTS idx_user_tenant_memberships_active;
DROP INDEX CONCURRENTLY IF EXISTS idx_user_tenant_memberships_user_active;

DROP INDEX CONCURRENTLY IF EXISTS idx_audit_logs_user_id;
DROP INDEX CONCURRENTLY IF EXISTS idx_audit_logs_tenant_id;
DROP INDEX CONCURRENTLY IF EXISTS idx_audit_logs_action;
DROP INDEX CONCURRENTLY IF EXISTS idx_audit_logs_created_at;
DROP INDEX CONCURRENTLY IF EXISTS idx_audit_logs_tenant_action;

DROP INDEX CONCURRENTLY IF EXISTS idx_notifications_user_id;
DROP INDEX CONCURRENTLY IF EXISTS idx_notifications_type;
DROP INDEX CONCURRENTLY IF EXISTS idx_notifications_read;
DROP INDEX CONCURRENTLY IF EXISTS idx_notifications_created_at;
DROP INDEX CONCURRENTLY IF EXISTS idx_notifications_user_created;

DROP INDEX CONCURRENTLY IF EXISTS idx_users_name_email;
DROP INDEX CONCURRENTLY IF EXISTS idx_feedback_search;

DROP INDEX CONCURRENTLY IF EXISTS idx_feedback_recent_unmoderated;
DROP INDEX CONCURRENTLY IF EXISTS idx_notifications_unread_recent;

-- Reset autovacuum settings to defaults
ALTER TABLE feedback RESET (autovacuum_vacuum_scale_factor);
ALTER TABLE feedback RESET (autovacuum_analyze_scale_factor);
ALTER TABLE comments RESET (autovacuum_vacuum_scale_factor);
ALTER TABLE reactions RESET (autovacuum_vacuum_scale_factor);
ALTER TABLE notifications RESET (autovacuum_vacuum_scale_factor);
