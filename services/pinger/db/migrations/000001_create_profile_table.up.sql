CREATE TABLE IF NOT EXISTS profile(
    id serial PRIMARY KEY,
    guid uuid UNIQUE NOT NULL,
    name VARCHAR (100) NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL
);