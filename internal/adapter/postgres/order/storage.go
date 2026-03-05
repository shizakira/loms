package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
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

func (s *Storage) Create(ctx context.Context, order domain.Order) (domain.Order, error) {
	q := s.query(ctx)

	created, err := q.CreateOrder(ctx, sqlc.CreateOrderParams{
		Status: sqlc.OrderStatus(order.Status),
		UserID: int64(order.User),
	})
	if err != nil {
		return domain.Order{}, fmt.Errorf("q.CreateOrder: %w", err)
	}

	for _, item := range order.Items {
		err = q.CreateOrderItem(ctx, sqlc.CreateOrderItemParams{
			OrderID: created.ID,
			Sku:     int64(item.Sku),
			Count:   int16(item.Count),
		})
		if err != nil {
			return domain.Order{}, fmt.Errorf("q.CreateOrderItem sku=%d: %w", item.Sku, err)
		}
	}

	order.ID = int(created.ID)
	return order, nil
}

func (s *Storage) GetByID(ctx context.Context, orderID int) (domain.Order, error) {
	q := s.query(ctx)

	row, err := q.GetOrderByID(ctx, int64(orderID))
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
			Sku:   uint32(item.Sku),
			Count: uint16(item.Count),
		})
	}

	return domain.Order{
		ID:     int(row.ID),
		Status: domain.OrderStatus(row.Status),
		User:   int(row.UserID),
		Items:  items,
	}, nil
}

func (s *Storage) Save(ctx context.Context, order domain.Order) error {
	q := s.query(ctx)

	err := q.SaveOrder(ctx, sqlc.SaveOrderParams{
		ID:     int64(order.ID),
		Status: sqlc.OrderStatus(order.Status),
	})
	if err != nil {
		return fmt.Errorf("q.SaveOrder: %w", err)
	}

	return nil
}
