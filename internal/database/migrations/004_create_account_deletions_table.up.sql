-- Create account_deletions table for scheduled account deletions
CREATE TABLE IF NOT EXISTS account_deletions (
    deletion_id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    scheduled_at TIMESTAMP WITH TIME ZONE NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index on user_id for faster lookups
CREATE INDEX IF NOT EXISTS idx_account_deletions_user_id ON account_deletions(user_id);

-- Create index on scheduled_at for cleanup queries
CREATE INDEX IF NOT EXISTS idx_account_deletions_scheduled_at ON account_deletions(scheduled_at);

