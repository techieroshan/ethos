-- Create moderation_appeals table for moderation appeals
CREATE TABLE IF NOT EXISTS moderation_appeals (
    appeal_id VARCHAR(255) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    moderated_item_id VARCHAR(255) NOT NULL,
    item_type VARCHAR(50) NOT NULL, -- 'feedback', 'comment', 'profile', etc.
    reason VARCHAR(500) NOT NULL,
    details TEXT,
    status VARCHAR(50) DEFAULT 'pending', -- 'pending', 'warned', 'actioned', 'escalated', 'appealed'
    submitted_by VARCHAR(255) NOT NULL REFERENCES users(id),
    submitted_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    reviewed_at TIMESTAMP WITH TIME ZONE,
    reviewer_notes TEXT,

    FOREIGN KEY (submitted_by) REFERENCES users(id)
);

-- Create indexes for efficient querying
CREATE INDEX IF NOT EXISTS idx_moderation_appeals_item ON moderation_appeals(moderated_item_id, item_type);
CREATE INDEX IF NOT EXISTS idx_moderation_appeals_submitted_by ON moderation_appeals(submitted_by);
CREATE INDEX IF NOT EXISTS idx_moderation_appeals_status ON moderation_appeals(status);
CREATE INDEX IF NOT EXISTS idx_moderation_appeals_submitted_at ON moderation_appeals(submitted_at);
