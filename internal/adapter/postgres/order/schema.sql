CREATE TYPE order_status AS ENUM (
    'new',
    'awaiting_payment',
    'failed',
    'payed',
    'cancelled'
    );

CREATE TABLE orders
(
    id      BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    status  order_status NOT NULL DEFAULT 'new',
    user_id BIGINT       NOT NULL
);

CREATE TABLE order_items
(
    order_id BIGINT   NOT NULL REFERENCES orders (id),
    sku      BIGINT   NOT NULL,
    count    SMALLINT NOT NULL CHECK ( count > 0),
    PRIMARY KEY (order_id, sku)
);