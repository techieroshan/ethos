-- Create feedback_reactions table
CREATE TABLE IF NOT EXISTS feedback_reactions (
    reaction_id VARCHAR(255) PRIMARY KEY,
    feedback_id VARCHAR(255) NOT NULL REFERENCES feedback_items(feedback_id) ON DELETE CASCADE,
    user_id VARCHAR(255) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    reaction_type VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(feedback_id, user_id, reaction_type)
);

-- Create index on feedback_id for faster lookups
CREATE INDEX IF NOT EXISTS idx_feedback_reactions_feedback_id ON feedback_reactions(feedback_id);

-- Create index on user_id for faster lookups
CREATE INDEX IF NOT EXISTS idx_feedback_reactions_user_id ON feedback_reactions(user_id);

-- Create index on reaction_type for aggregation
CREATE INDEX IF NOT EXISTS idx_feedback_reactions_reaction_type ON feedback_reactions(reaction_type);

