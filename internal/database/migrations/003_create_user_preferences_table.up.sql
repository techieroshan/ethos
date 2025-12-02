-- Create user_preferences table
CREATE TABLE IF NOT EXISTS user_preferences (
    user_id VARCHAR(255) PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    notify_on_login BOOLEAN DEFAULT TRUE,
    locale VARCHAR(10) DEFAULT 'en-US',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index on user_id (already primary key, but explicit for clarity)
CREATE INDEX IF NOT EXISTS idx_user_preferences_user_id ON user_preferences(user_id);

