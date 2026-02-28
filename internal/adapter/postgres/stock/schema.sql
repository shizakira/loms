CREATE TABLE stocks
(
    sku_id      BIGINT NOT NULL PRIMARY KEY,
    total_count BIGINT NOT NULL DEFAULT 0,
    reserved    BIGINT NOT NULL DEFAULT 0
);