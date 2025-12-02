-- Create notification_preferences table
CREATE TABLE IF NOT EXISTS notification_preferences (
    user_id VARCHAR(255) PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    email BOOLEAN DEFAULT TRUE,
    push BOOLEAN DEFAULT TRUE,
    in_app BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index on user_id (already primary key, but explicit for clarity)
CREATE INDEX IF NOT EXISTS idx_notification_preferences_user_id ON notification_preferences(user_id);

