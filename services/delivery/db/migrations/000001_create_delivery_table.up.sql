CREATE TABLE IF NOT EXISTS delivery(
    id serial PRIMARY KEY,
    order_id uuid NOT NULL,
    user_id uuid NOT NULL,
    status integer NOT NULL,
    planned_date_start timestamptz NOT NULL,
    planned_date_end timestamptz NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL
);