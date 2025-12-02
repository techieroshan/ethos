-- Create feedback_follow_ups table for follow-up discussions
CREATE TABLE IF NOT EXISTS feedback_follow_ups (
    follow_up_id VARCHAR(255) PRIMARY KEY,
    feedback_id VARCHAR(255) NOT NULL REFERENCES feedback_items(feedback_id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    author_id VARCHAR(255) NOT NULL REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (author_id) REFERENCES users(id)
);

-- Create index for efficient querying by feedback_id
CREATE INDEX IF NOT EXISTS idx_feedback_follow_ups_feedback_id ON feedback_follow_ups(feedback_id);
CREATE INDEX IF NOT EXISTS idx_feedback_follow_ups_created_at ON feedback_follow_ups(created_at);
