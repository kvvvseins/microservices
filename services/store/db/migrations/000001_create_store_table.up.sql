CREATE TABLE IF NOT EXISTS store(
    id serial PRIMARY KEY,
    guid uuid UNIQUE NOT NULL,
    name varchar(255) UNIQUE NOT NULL,
    price integer NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL
);