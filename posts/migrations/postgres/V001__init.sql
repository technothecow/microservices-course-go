CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create the posts table
CREATE TABLE posts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    creator_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_private BOOLEAN NOT NULL DEFAULT FALSE,
    tags TEXT[] NOT NULL DEFAULT '{}'
);

CREATE INDEX idx_posts_creator_id ON posts (creator_id);
CREATE INDEX idx_posts_created_at ON posts (created_at);
CREATE INDEX idx_posts_updated_at ON posts (updated_at);
CREATE INDEX idx_posts_is_private ON posts (is_private);
-- GIN (Generalized Inverted Index) is optimized for handling arrays and full-text search
-- It creates an index entry for each element in the array, making array operations like 
-- containment and overlap very efficient
CREATE INDEX idx_posts_tags ON posts USING GIN (tags);

CREATE OR REPLACE FUNCTION update_posts_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_posts_updated_at
BEFORE UPDATE ON posts
FOR EACH ROW
EXECUTE FUNCTION update_posts_updated_at();