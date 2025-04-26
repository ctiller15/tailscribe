-- +goose Up
CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    -- Validate email in-code before sending to DB.
    email TEXT UNIQUE,
    username TEXT UNIQUE,
    firstname TEXT,
    lastname TEXT,
    -- Validate password emptiness before sending to DB.
    password TEXT,
    facebook_id TEXT,
    reset_password_token TEXT,
    reset_password_expires TEXT,
    -- Deprecated. To remove premium version.
    is_premium BOOLEAN NOT NULL DEFAULT FALSE,
    premium_level INT NOT NULL DEFAULT 0,
    stripe_customer_id TEXT,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    created_at DATE NOT NULL,
    updated_at DATE NOT NULL
);

-- +goose Down
DROP TABLE users;