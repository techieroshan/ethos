-- Create notifications table
CREATE TABLE IF NOT EXISTS notifications (
    notification_id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    message TEXT NOT NULL,
    read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index on user_id for faster lookups
CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);

-- Create index on read status for filtering
CREATE INDEX IF NOT EXISTS idx_notifications_read ON notifications(read);

-- Create index on created_at for sorting
CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON notifications(created_at DESC);

