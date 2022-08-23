BEGIN;

CREATE TABLE IF NOT EXISTS withdrawals (
    id SERIAL PRIMARY KEY,
    user_id SERIAL REFERENCES users(id) ON DELETE CASCADE,
    order_number VARCHAR (255) NOT NULL UNIQUE,
    processed_at TIMESTAMP,
    sum FLOAT DEFAULT 0
);

CREATE INDEX withdrawals_processed_at_id_idx ON withdrawals (processed_at, id);

COMMIT;