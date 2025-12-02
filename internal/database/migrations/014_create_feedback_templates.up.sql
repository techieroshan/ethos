-- Create feedback_templates table for template management
CREATE TABLE IF NOT EXISTS feedback_templates (
    template_id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    context_tags TEXT[], -- Array of context tags
    template_fields JSONB, -- JSON schema for template fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index on context_tags for efficient filtering
CREATE INDEX IF NOT EXISTS idx_feedback_templates_context_tags ON feedback_templates USING GIN(context_tags);
