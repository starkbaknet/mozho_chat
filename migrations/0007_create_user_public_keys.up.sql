CREATE TABLE user_public_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    public_key TEXT NOT NULL,
    key_version VARCHAR(50) DEFAULT 'v1',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    expires_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_user_public_keys_user_id ON user_public_keys(user_id);
