package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/shizakira/loms/internal/adapter/postgres/order/sqlc"
	"github.com/shizakira/loms/internal/domain"
	"github.com/shizakira/loms/pkg/transaction"
)

type Storage struct{}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) query(ctx context.Context) *sqlc.Queries {
	db := transaction.TryExtractTX(ctx)

	return sqlc.New(db)
}

func (s *Storage) Create(ctx context.Context, order domain.Order) (int, error) {
	q := s.query(ctx)

	orderID, err := q.CreateOrder(ctx, sqlc.CreateOrderParams{
		Status: sqlc.OrderStatus(order.Status),
		UserID: int64(order.User),
	})
	if err != nil {
		return 0, fmt.Errorf("q.CreateOrder: %w", err)
	}

	for _, item := range order.Items {
		err = q.CreateOrderItem(ctx, sqlc.CreateOrderItemParams{
			OrderID: orderID,
			Sku:     int64(item.Sku),
			Count:   int16(item.Count),
		})
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23514" {
				return 0, domain.ErrInsufficientStock
			}

			return 0, fmt.Errorf("q.CreateOrderItem sku=%d: %w", item.Sku, err)
		}
	}

	return int(orderID), nil
}

func (s *Storage) GetByID(ctx context.Context, orderID int, pessimistic bool) (domain.Order, error) {
	q := s.query(ctx)

	var row sqlc.Order
	var err error
	if pessimistic {
		row, err = q.GetOrderByIDForUpdate(ctx, int64(orderID))
	} else {
		row, err = q.GetOrderByID(ctx, int64(orderID))
	}
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Order{}, domain.ErrOrderNotFound
		}
		return domain.Order{}, fmt.Errorf("q.GetOrderByID: %w", err)
	}

	itemRows, err := q.GetOrderItems(ctx, row.ID)
	if err != nil {
		return domain.Order{}, fmt.Errorf("q.GetOrderItems: %w", err)
	}

	items := make([]domain.OrderItem, 0, len(itemRows))
	for _, item := range itemRows {
		items = append(items, domain.OrderItem{
			Sku:   int(item.Sku),
			Count: int(item.Count),
		})
	}

	return domain.Order{
		ID:     int(row.ID),
		Status: domain.OrderStatus(row.Status),
		User:   int(row.UserID),
		Items:  items,
	}, nil
}

func (s *Storage) UpdateStatus(ctx context.Context, orderID int, status domain.OrderStatus) error {
	q := s.query(ctx)

	err := q.UpdateOrderStatus(ctx, sqlc.UpdateOrderStatusParams{
		ID:     int64(orderID),
		Status: sqlc.OrderStatus(status),
	})
	if err != nil {
		return fmt.Errorf("q.UpdateStatus: %w", err)
	}

	return nil
}
