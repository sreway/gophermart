BEGIN;

CREATE TABLE IF NOT EXISTS balance(
    id uuid PRIMARY KEY,
    user_id uuid REFERENCES users(id) ON DELETE CASCADE,
    balance float NOT NULL DEFAULT 0,
    withdrawn float NOT NULL DEFAULT 0
);

COMMIT;