-- name: GetStockBySku :one
SELECT sku_id, total_count, reserved
FROM stocks
WHERE sku_id = @sku_id;

-- name: DecreaseReservedStock :exec
UPDATE stocks
SET reserved = reserved - @count
WHERE sku_id = @sku_id;

-- name: DecreaseReserveAndTotalCountStock :exec
UPDATE stocks
SET reserved    = reserved - @count,
    total_count = total_count - @count
WHERE sku_id = @sku_id;

-- name: IncreaseReservedStock :execrows
UPDATE stocks
SET reserved = reserved + @count
WHERE sku_id = @sku_id;