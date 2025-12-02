-- Create feedback_items table
CREATE TABLE IF NOT EXISTS feedback_items (
    feedback_id VARCHAR(255) PRIMARY KEY,
    author_id VARCHAR(255) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    type VARCHAR(50),
    visibility VARCHAR(50) DEFAULT 'public',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index on author_id for faster lookups
CREATE INDEX IF NOT EXISTS idx_feedback_items_author_id ON feedback_items(author_id);

-- Create index on created_at for sorting
CREATE INDEX IF NOT EXISTS idx_feedback_items_created_at ON feedback_items(created_at DESC);

-- Create index on visibility for filtering
CREATE INDEX IF NOT EXISTS idx_feedback_items_visibility ON feedback_items(visibility);

-- Create index on type for filtering
CREATE INDEX IF NOT EXISTS idx_feedback_items_type ON feedback_items(type);

