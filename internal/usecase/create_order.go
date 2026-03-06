package usecase

import (
	"context"
	"fmt"

	"github.com/shizakira/loms/internal/domain"
	"github.com/shizakira/loms/internal/dto"
	"github.com/shizakira/loms/pkg/transaction"
)

func (l *Loms) CreateOrder(ctx context.Context, input dto.CreateOrderInput) (dto.CreateOrderOutput, error) {
	items := l.aggregateItems(input.Items)

	newOrder := domain.Order{
		Status: domain.StatusNew,
		User:   input.User,
		Items:  items,
	}

	orderID, err := l.orderStorage.Create(ctx, newOrder)
	if err != nil {
		return dto.CreateOrderOutput{}, fmt.Errorf("orderStorage.Create: %w", err)
	}

	err = transaction.Wrap(ctx, func(ctx context.Context) error {
		order, err := l.orderStorage.GetByID(ctx, orderID, true)
		if err != nil {
			return fmt.Errorf("orderStorage.GetByID: %w", err)
		}

		for _, item := range order.Items {
			if err := l.stockStorage.Reserve(ctx, item.Sku, item.Count); err != nil {
				if err := l.orderStorage.UpdateStatus(ctx, orderID, domain.StatusFailed); err != nil {
					return fmt.Errorf("orderStorage.Save: %w", err)
				}

				return fmt.Errorf("stockStorage.Reserve: %w", err)
			}
		}

		if err := l.orderStorage.UpdateStatus(ctx, orderID, domain.StatusAwaitingPayment); err != nil {
			return fmt.Errorf("orderStorage.Save: %w", err)
		}

		return nil
	})

	if err != nil {
		return dto.CreateOrderOutput{}, fmt.Errorf("transaction.Wrap: %w", err)
	}

	return dto.CreateOrderOutput{OrderID: orderID}, nil
}

func (l *Loms) aggregateItems(items []dto.OrderItem) []domain.OrderItem {
	m := make(map[int]int, len(items))

	for _, it := range items {
		m[it.Sku] += it.Count
	}

	res := make([]domain.OrderItem, 0, len(m))
	for sku, cnt := range m {
		res = append(res, domain.OrderItem{
			Sku:   sku,
			Count: cnt,
		})
	}

	return res
}
