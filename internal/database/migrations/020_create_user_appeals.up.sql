-- Create user_appeals table for user-submitted appeals
CREATE TABLE IF NOT EXISTS user_appeals (
    appeal_id VARCHAR(255) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    user_id VARCHAR(255) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL, -- 'content_removal', 'account_suspension', 'feedback_removal', 'rating_dispute', 'other'
    reference_id VARCHAR(255), -- optional reference to the item being appealed
    description TEXT NOT NULL,
    status VARCHAR(50) DEFAULT 'pending', -- 'pending', 'under_review', 'approved', 'rejected', 'closed'
    admin_notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    resolved_at TIMESTAMP WITH TIME ZONE,

    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Create indexes for efficient querying
CREATE INDEX IF NOT EXISTS idx_user_appeals_user_id ON user_appeals(user_id);
CREATE INDEX IF NOT EXISTS idx_user_appeals_status ON user_appeals(status);
CREATE INDEX IF NOT EXISTS idx_user_appeals_created_at ON user_appeals(created_at);
CREATE INDEX IF NOT EXISTS idx_user_appeals_type ON user_appeals(type);
