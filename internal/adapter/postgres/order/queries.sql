-- name: CreateOrder :one
INSERT INTO orders (status, user_id)
VALUES (@status, @user_id)
    RETURNING id;

-- name: GetOrderByID :one
SELECT id, status, user_id
FROM orders
WHERE id = @id;

-- name: GetOrderByIDForUpdate :one
SELECT id, status, user_id
FROM orders
WHERE id = @id
FOR NO KEY UPDATE;

-- name: UpdateOrderStatus :exec
UPDATE orders
SET status = @status
WHERE id = @id;

-- name: CreateOrderItem :exec
INSERT INTO order_items (order_id, sku, count)
VALUES (@order_id, @sku, @count);

-- name: GetOrderItems :many
SELECT sku, count
FROM order_items
WHERE order_id = @order_id;