CREATE TABLE IF NOT EXISTS accounts (
    user_id BIGINT PRIMARY KEY,
    first_name TEXT,
    last_name TEXT,
    username TEXT,
    phone TEXT,
    about TEXT,
    birthday TEXT,
    personal_channel_id TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT (now())
);