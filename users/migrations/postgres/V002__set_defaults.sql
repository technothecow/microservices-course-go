-- Set default fields for users table (not reverse-compatible!)
ALTER TABLE users
    ALTER COLUMN full_name SET DEFAULT '',
    ALTER COLUMN phone_number SET DEFAULT '',
    ALTER COLUMN last_login SET DEFAULT CURRENT_TIMESTAMP;
