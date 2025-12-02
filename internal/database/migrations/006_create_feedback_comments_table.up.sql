-- Create feedback_comments table
CREATE TABLE IF NOT EXISTS feedback_comments (
    comment_id VARCHAR(255) PRIMARY KEY,
    feedback_id VARCHAR(255) NOT NULL REFERENCES feedback_items(feedback_id) ON DELETE CASCADE,
    author_id VARCHAR(255) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    parent_comment_id VARCHAR(255) REFERENCES feedback_comments(comment_id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index on feedback_id for faster lookups
CREATE INDEX IF NOT EXISTS idx_feedback_comments_feedback_id ON feedback_comments(feedback_id);

-- Create index on author_id for faster lookups
CREATE INDEX IF NOT EXISTS idx_feedback_comments_author_id ON feedback_comments(author_id);

-- Create index on parent_comment_id for threading
CREATE INDEX IF NOT EXISTS idx_feedback_comments_parent_comment_id ON feedback_comments(parent_comment_id);

-- Create index on created_at for sorting
CREATE INDEX IF NOT EXISTS idx_feedback_comments_created_at ON feedback_comments(created_at DESC);

