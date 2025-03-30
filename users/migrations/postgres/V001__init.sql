-- Initial migration for user management system
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table
CREATE TABLE users
(
    id            UUID PRIMARY KEY      DEFAULT uuid_generate_v4(),
    username      VARCHAR(50)  NOT NULL UNIQUE,
    email         VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    full_name     VARCHAR(100),
    phone_number  VARCHAR(20),
    created_at    TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_login    TIMESTAMP,
    is_active     BOOLEAN               DEFAULT TRUE
);

-- User profiles table with 1-to-1 relationship to users
CREATE TABLE user_profiles
(
    user_id             UUID PRIMARY KEY,
    bio                 TEXT,
    profile_picture_url VARCHAR(255),
    location            VARCHAR(100),
    website             VARCHAR(255),
    birth_date          DATE,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- Followers table for following relationships
CREATE TABLE followers
(
    id           UUID PRIMARY KEY   DEFAULT uuid_generate_v4(),
    follower_id  UUID      NOT NULL,
    following_id UUID      NOT NULL,
    created_at   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (follower_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (following_id) REFERENCES users (id) ON DELETE CASCADE,
    -- Prevent duplicate follows
    UNIQUE (follower_id, following_id)
);

-- Indexes for better performance
CREATE INDEX idx_users_username ON users (username);
CREATE INDEX idx_users_email ON users (email);
CREATE INDEX idx_users_phone_number ON users(phone_number);
CREATE INDEX idx_followers_follower ON followers (follower_id);
CREATE INDEX idx_followers_following ON followers (following_id);

-- Create trigger function to automatically update the updated_at column
CREATE OR REPLACE FUNCTION update_users_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';

-- Create trigger to automatically update updated_at column on each row update
CREATE TRIGGER trigger_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_users_updated_at();
