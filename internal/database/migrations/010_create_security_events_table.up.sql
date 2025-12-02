-- Create security_events table
CREATE TABLE IF NOT EXISTS security_events (
    event_id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    ip VARCHAR(45),
    location VARCHAR(100),
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index on user_id for faster lookups
CREATE INDEX IF NOT EXISTS idx_security_events_user_id ON security_events(user_id);

-- Create index on timestamp for sorting
CREATE INDEX IF NOT EXISTS idx_security_events_timestamp ON security_events(timestamp DESC);

