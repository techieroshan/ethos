-- Performance optimization indexes and constraints

-- Users table indexes
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_created_at ON users(created_at);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_email_verified ON users(email_verified) WHERE email_verified = false;

-- Feedback table indexes
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_feedback_author_id ON feedback(author_id);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_feedback_recipient_id ON feedback(recipient_id);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_feedback_created_at ON feedback(created_at DESC);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_feedback_visibility ON feedback(visibility);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_feedback_rating ON feedback(rating);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_feedback_tenant_author ON feedback(tenant_id, author_id);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_feedback_tenant_created ON feedback(tenant_id, created_at DESC);

-- Full-text search index for feedback content
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_feedback_content_fts ON feedback USING gin(to_tsvector('english', content));

-- Comments table indexes
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_comments_feedback_id ON comments(feedback_id);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_comments_author_id ON comments(author_id);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_comments_created_at ON comments(created_at DESC);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_comments_tenant_feedback ON comments(tenant_id, feedback_id);

-- Reactions table indexes
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_reactions_feedback_id ON reactions(feedback_id);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_reactions_user_id ON reactions(user_id);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_reactions_type ON reactions(reaction_type);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_reactions_tenant_user ON reactions(tenant_id, user_id);

-- User role assignments indexes
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_user_role_assignments_user_id ON user_role_assignments(user_id);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_user_role_assignments_role_id ON user_role_assignments(role_id);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_user_role_assignments_active ON user_role_assignments(is_active) WHERE is_active = true;
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_user_role_assignments_user_active ON user_role_assignments(user_id, is_active) WHERE is_active = true;

-- User tenant memberships indexes
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_user_tenant_memberships_user_id ON user_tenant_memberships(user_id);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_user_tenant_memberships_tenant_id ON user_tenant_memberships(tenant_id);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_user_tenant_memberships_active ON user_tenant_memberships(is_active) WHERE is_active = true;
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_user_tenant_memberships_user_active ON user_tenant_memberships(user_id, is_active) WHERE is_active = true;

-- Audit logs indexes (for performance)
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_audit_logs_tenant_id ON audit_logs(tenant_id);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_audit_logs_action ON audit_logs(action);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_audit_logs_created_at ON audit_logs(created_at DESC);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_audit_logs_tenant_action ON audit_logs(tenant_id, action);

-- Notifications indexes
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_notifications_type ON notifications(notification_type);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_notifications_read ON notifications(is_read) WHERE is_read = false;
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_notifications_created_at ON notifications(created_at DESC);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_notifications_user_created ON notifications(user_id, created_at DESC);

-- Search optimization indexes
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_name_email ON users USING gin(to_tsvector('english', name || ' ' || COALESCE(public_bio, '')));
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_feedback_search ON feedback USING gin(to_tsvector('english', content || ' ' || COALESCE(title, '')));

-- Partitioning for large tables (if needed in future)
-- These are commented out but ready to implement when tables grow large

-- CREATE TABLE feedback_2024 PARTITION OF feedback FOR VALUES FROM ('2024-01-01') TO ('2025-01-01');
-- CREATE TABLE audit_logs_2024 PARTITION OF audit_logs FOR VALUES FROM ('2024-01-01') TO ('2025-01-01');

-- Optimize autovacuum settings for high-traffic tables
ALTER TABLE feedback SET (autovacuum_vacuum_scale_factor = 0.02);
ALTER TABLE feedback SET (autovacuum_analyze_scale_factor = 0.01);
ALTER TABLE comments SET (autovacuum_vacuum_scale_factor = 0.05);
ALTER TABLE reactions SET (autovacuum_vacuum_scale_factor = 0.05);
ALTER TABLE notifications SET (autovacuum_vacuum_scale_factor = 0.1);

-- Create partial indexes for common queries
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_feedback_recent_unmoderated ON feedback(created_at DESC, visibility)
WHERE visibility = 'pending_moderation' AND created_at > NOW() - INTERVAL '30 days';

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_notifications_unread_recent ON notifications(user_id, created_at DESC)
WHERE is_read = false AND created_at > NOW() - INTERVAL '7 days';
