-- Create user sessions and organization context tracking
-- +migrate Up

-- Create user sessions table
CREATE TABLE IF NOT EXISTS user_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(255) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    organization_id UUID REFERENCES organizations(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) UNIQUE NOT NULL,
    refresh_token_hash VARCHAR(255) UNIQUE,
    ip_address INET,
    user_agent TEXT,
    device_name VARCHAR(255),
    last_activity_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    revoked_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create organization activity log
CREATE TABLE IF NOT EXISTS organization_activity_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    user_id VARCHAR(255) REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(50) NOT NULL, -- created, updated, deleted, invited, removed, etc.
    resource_type VARCHAR(50) NOT NULL, -- user, feedback, settings, etc.
    resource_id VARCHAR(255),
    changes JSONB DEFAULT '{}',
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create user context switches log (for audit trail)
CREATE TABLE IF NOT EXISTS user_context_switches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(255) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    from_organization_id UUID REFERENCES organizations(id) ON DELETE SET NULL,
    to_organization_id UUID REFERENCES organizations(id) ON DELETE SET NULL,
    session_id UUID NOT NULL REFERENCES user_sessions(id) ON DELETE CASCADE,
    ip_address INET,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Add indexes for performance
CREATE INDEX IF NOT EXISTS idx_user_sessions_user_id ON user_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_user_sessions_organization_id ON user_sessions(organization_id);
CREATE INDEX IF NOT EXISTS idx_user_sessions_expires_at ON user_sessions(expires_at);
CREATE INDEX IF NOT EXISTS idx_user_sessions_revoked_at ON user_sessions(revoked_at);
CREATE INDEX IF NOT EXISTS idx_organization_activity_log_org_id ON organization_activity_log(organization_id);
CREATE INDEX IF NOT EXISTS idx_organization_activity_log_user_id ON organization_activity_log(user_id);
CREATE INDEX IF NOT EXISTS idx_organization_activity_log_created_at ON organization_activity_log(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_organization_activity_log_action ON organization_activity_log(action);
CREATE INDEX IF NOT EXISTS idx_user_context_switches_user_id ON user_context_switches(user_id);
CREATE INDEX IF NOT EXISTS idx_user_context_switches_timestamp ON user_context_switches(timestamp DESC);
