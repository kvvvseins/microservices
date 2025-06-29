CREATE TABLE IF NOT EXISTS reserve(
    id serial PRIMARY KEY,
    store_id int NOT NULL,
    quantity int NOT NULL,
    order_id uuid NULL DEFAULT NULL,
    user_id uuid NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL
);

alter table reserve
    add constraint reserve_store_id_fk
        foreign key (store_id) references store
            on update restrict on delete cascade;
