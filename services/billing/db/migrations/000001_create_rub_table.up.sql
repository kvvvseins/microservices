CREATE TABLE IF NOT EXISTS rub(
    id serial PRIMARY KEY,
    user_id uuid UNIQUE NOT NULL,
    value int NOT NULL default 0 CHECK (value >= 0),
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL
);