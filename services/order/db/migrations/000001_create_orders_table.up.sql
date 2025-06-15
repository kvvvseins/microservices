CREATE TABLE IF NOT EXISTS orders(
    id serial PRIMARY KEY,
    guid uuid UNIQUE NOT NULL,
    user_id uuid NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL
);