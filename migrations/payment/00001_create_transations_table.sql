-- +goose Up
CREATE TABLE transactions(
    uuid UUID primary key ,
    order_uuid UUID NOT NULL,
    status varchar(50) not null default 'PROCESSING',
    created_at timestamp not null default NOW(),
    updated_at timestamp
);

-- +goose Down
DROP TABLE IF EXISTS transactions;
