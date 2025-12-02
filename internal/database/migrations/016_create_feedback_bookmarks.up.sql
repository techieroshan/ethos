-- Create feedback_bookmarks table for user bookmarks
CREATE TABLE IF NOT EXISTS feedback_bookmarks (
    bookmark_id VARCHAR(255) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    user_id VARCHAR(255) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    feedback_id VARCHAR(255) NOT NULL REFERENCES feedback_items(feedback_id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(user_id, feedback_id), -- Prevent duplicate bookmarks

    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (feedback_id) REFERENCES feedback_items(feedback_id)
);

-- Create indexes for efficient querying
CREATE INDEX IF NOT EXISTS idx_feedback_bookmarks_user_id ON feedback_bookmarks(user_id);
CREATE INDEX IF NOT EXISTS idx_feedback_bookmarks_feedback_id ON feedback_bookmarks(feedback_id);
CREATE INDEX IF NOT EXISTS idx_feedback_bookmarks_created_at ON feedback_bookmarks(created_at);
