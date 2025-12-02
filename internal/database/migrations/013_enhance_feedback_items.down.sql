-- Remove enhanced columns from feedback_items table
ALTER TABLE feedback_items
DROP COLUMN IF EXISTS is_anonymous,
DROP COLUMN IF EXISTS helpfulness,
DROP COLUMN IF EXISTS reviewer_context,
DROP COLUMN IF EXISTS moderation_state;
