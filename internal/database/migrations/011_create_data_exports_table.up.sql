-- Create data_exports table
CREATE TABLE IF NOT EXISTS data_exports (
    export_id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(50) DEFAULT 'pending',
    download_url TEXT,
    expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index on user_id for faster lookups
CREATE INDEX IF NOT EXISTS idx_data_exports_user_id ON data_exports(user_id);

-- Create index on status for filtering
CREATE INDEX IF NOT EXISTS idx_data_exports_status ON data_exports(status);

