BEGIN;

CREATE TABLE IF NOT EXISTS orders(
    id uuid PRIMARY KEY,
    number bigint NOT NULL UNIQUE,
    user_id uuid REFERENCES users(id) ON DELETE CASCADE,
    status OrderStatus NOT NULL,
    accrual float NOT NULL DEFAULT 0,
    uploaded_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMIT;