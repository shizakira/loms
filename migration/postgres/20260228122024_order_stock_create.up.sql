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
    count    SMALLINT NOT NULL,
    PRIMARY KEY (order_id, sku)
);

CREATE TABLE stocks
(
    sku_id      BIGINT NOT NULL PRIMARY KEY,
    total_count BIGINT NOT NULL DEFAULT 0,
    reserved    BIGINT NOT NULL DEFAULT 0
);

INSERT INTO stocks (sku_id, total_count, reserved)
VALUES (1076963, 100, 0),
       (1148162, 50, 0),
       (1625903, 200, 0);