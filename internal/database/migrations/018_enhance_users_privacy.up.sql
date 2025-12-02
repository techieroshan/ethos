-- Add privacy control columns to users table
ALTER TABLE users
ADD COLUMN IF NOT EXISTS opt_outs TEXT[], -- Array of opt-out preferences
ADD COLUMN IF NOT EXISTS anonymized_at TIMESTAMP WITH TIME ZONE,
ADD COLUMN IF NOT EXISTS delete_requested_at TIMESTAMP WITH TIME ZONE;

-- Create index on opt_outs for efficient querying
CREATE INDEX IF NOT EXISTS idx_users_opt_outs ON users USING GIN(opt_outs);
