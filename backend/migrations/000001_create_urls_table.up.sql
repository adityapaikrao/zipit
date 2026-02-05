CREATE TABLE IF NOT EXISTS urls (
    id BIGSERIAL PRIMARY KEY,
    long_url TEXT NOT NULL,
    short_code VARCHAR(12) UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Add an index on long_url for fast existence checks during shortening
CREATE INDEX IF NOT EXISTS idx_urls_long_url ON urls(long_url);
