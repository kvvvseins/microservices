CREATE TABLE IF NOT EXISTS notify(
    id serial PRIMARY KEY,
    user_id uuid UNIQUE NOT NULL,
    email varchar(256) NOT NULL,
    message text NOT NULL,
    value int NOT NULL default 0 CHECK (value >= 0),
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL
);