-- Add user roles system for admin functionality
-- +migrate Up

-- Create user roles table
CREATE TABLE IF NOT EXISTS user_roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    permissions JSONB DEFAULT '{}',
    is_system_role BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create user role assignments table
CREATE TABLE IF NOT EXISTS user_role_assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(255) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES user_roles(id) ON DELETE CASCADE,
    assigned_by VARCHAR(255) REFERENCES users(id),
    assigned_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE,
    is_active BOOLEAN DEFAULT true,
    UNIQUE(user_id, role_id)
);

-- Add indexes for performance
CREATE INDEX IF NOT EXISTS idx_user_role_assignments_user_id ON user_role_assignments(user_id);
CREATE INDEX IF NOT EXISTS idx_user_role_assignments_role_id ON user_role_assignments(role_id);
CREATE INDEX IF NOT EXISTS idx_user_role_assignments_active ON user_role_assignments(is_active) WHERE is_active = true;

-- Insert default roles
INSERT INTO user_roles (name, description, permissions, is_system_role) VALUES
('user', 'Standard user with basic access', '{"read_feedback": true, "write_feedback": true, "read_profile": true, "write_profile": true}', true),
('moderator', 'Content moderator with review capabilities', '{"read_feedback": true, "moderate_content": true, "read_reports": true, "manage_appeals": true}', true),
('org_admin', 'Organization administrator', '{"manage_users": true, "manage_content": true, "view_analytics": true, "manage_settings": true}', true),
('platform_admin', 'Platform administrator with full access', '{"manage_all_users": true, "manage_all_content": true, "system_admin": true, "view_all_analytics": true, "manage_system_settings": true}', true);

-- Add role to existing users (default to 'user' role)
INSERT INTO user_role_assignments (user_id, role_id, assigned_by)
SELECT u.id, r.id, NULL
FROM users u
CROSS JOIN (SELECT id FROM user_roles WHERE name = 'user') r;

-- Add updated_at trigger for user_roles
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_user_roles_updated_at BEFORE UPDATE ON user_roles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
