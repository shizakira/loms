-- name: GetStockBySku :one
SELECT sku_id, total_count, reserved
FROM stocks
WHERE sku_id = @sku_id;

-- name: GetStocksBySkus :many
SELECT sku_id, total_count, reserved
FROM stocks
WHERE sku_id = ANY(@sku_ids::bigint[]);

-- name: SaveStock :exec
UPDATE stocks
SET total_count = @total_count,
    reserved    = @reserved
WHERE sku_id = @sku_id;