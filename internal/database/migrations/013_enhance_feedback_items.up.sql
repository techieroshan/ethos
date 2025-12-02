-- Add new columns to feedback_items table for enhanced features
ALTER TABLE feedback_items
ADD COLUMN IF NOT EXISTS is_anonymous BOOLEAN DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS helpfulness DOUBLE PRECISION DEFAULT 0.0,
ADD COLUMN IF NOT EXISTS reviewer_context JSONB,
ADD COLUMN IF NOT EXISTS moderation_state VARCHAR(50) DEFAULT 'pending';
