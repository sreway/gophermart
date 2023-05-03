BEGIN;

CREATE TABLE IF NOT EXISTS withdrawals (
    id uuid PRIMARY KEY,
    user_id uuid REFERENCES users(id) ON DELETE CASCADE,
    order_number bigint NOT NULL UNIQUE,
    processed_at TIMESTAMP,
    sum float DEFAULT 0
);

CREATE INDEX withdrawals_processed_at_id_idx ON withdrawals (processed_at, id);

COMMIT;