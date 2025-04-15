CREATE TABLE IF NOT EXISTS users(
    id serial PRIMARY KEY,
    guid uuid UNIQUE NOT NULL,
    password VARCHAR (100) NOT NULL,
    email VARCHAR (300) UNIQUE NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL
);