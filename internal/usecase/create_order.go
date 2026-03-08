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
	order := domain.Order{
		Status: domain.StatusNew,
		User:   input.User,
		Items:  items,
	}

	err := transaction.Wrap(ctx, func(ctx context.Context) error {
		var err error
		order.ID, err = l.orderStorage.Create(ctx, order)
		if err != nil {
			return fmt.Errorf("orderStorage.Create: %w", err)
		}

		ev, err := domain.NewOrderEvent(order.ID, domain.StatusNew)
		if err != nil {
			return fmt.Errorf("domain.NewOrderEvent: %w", err)
		}

		if err = l.outboxStorage.CreateEvent(ctx, ev); err != nil {
			return fmt.Errorf("outboxStorage.CreateEvent: %w", err)
		}

		return nil
	})
	if err != nil {
		return dto.CreateOrderOutput{}, fmt.Errorf("transaction.Wrap: %w", err)
	}

	err = transaction.Wrap(ctx, func(ctx context.Context) error {
		newStatus := domain.StatusAwaitingPayment

		reserveErr := transaction.WithSavepoint(ctx, "reserve", func(ctx context.Context) error {
			for _, item := range order.Items {
				if err := l.stockStorage.Reserve(ctx, item.Sku, item.Count); err != nil {
					return err
				}
			}
			return nil
		})
		if reserveErr != nil {
			newStatus = domain.StatusFailed
		}

		if err := l.orderStorage.UpdateStatus(ctx, order.ID, newStatus); err != nil {
			return fmt.Errorf("orderStorage.UpdateStatus: %w", err)
		}

		ev, err := domain.NewOrderEvent(order.ID, newStatus, domain.OrderEventWithError(reserveErr))
		if err != nil {
			return fmt.Errorf("domain.NewOrderEvent: %w", err)
		}

		if err = l.outboxStorage.CreateEvent(ctx, ev); err != nil {
			return fmt.Errorf("outboxStorage.CreateEvent: %w", err)
		}

		if reserveErr != nil {
			return reserveErr
		}

		return nil
	})
	if err != nil {
		return dto.CreateOrderOutput{}, fmt.Errorf("transaction.Wrap: %w", err)
	}

	return dto.CreateOrderOutput{OrderID: order.ID}, nil
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
