BEGIN;
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL
);
CREATE TABLE sessions (
    session_id TEXT PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    expires_at INTEGER NOT NULL
);
ALTER TABLE expenses ADD COLUMN user_id INTEGER REFERENCES users(id);
PRAGMA user_version = 3;
COMMIT;
