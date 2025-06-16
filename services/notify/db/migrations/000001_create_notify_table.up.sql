CREATE TABLE IF NOT EXISTS notify(
    id serial PRIMARY KEY,
    user_id uuid NOT NULL,
    email varchar(256) NOT NULL,
    message text NOT NULL,
    status smallint NOT NULL DEFAULT 1,
    value int NOT NULL default 0 CHECK (value >= 0),
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL
);