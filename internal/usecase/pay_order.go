package usecase

import (
	"context"
	"fmt"

	"github.com/shizakira/loms/internal/domain"
	"github.com/shizakira/loms/internal/dto"
	"github.com/shizakira/loms/pkg/transaction"
)

func (l *Loms) PayOrder(ctx context.Context, input dto.PayOrderInput) error {
	err := transaction.Wrap(ctx, func(ctx context.Context) error {
		order, err := l.orderStorage.GetByID(ctx, input.OrderID, true)
		if err != nil {
			return fmt.Errorf("orderStorage.GetByID: %w", err)
		}

		if order.Status != domain.StatusAwaitingPayment {
			return domain.ErrInvalidOrderStatus
		}

		for _, item := range order.Items {
			if err := l.stockStorage.DecreaseReserveAndTotalCount(ctx, item.Sku, item.Count); err != nil {
				return fmt.Errorf("stockStorage.DecreaseReserveAndTotalCount: %w", err)
			}
		}

		if err := l.orderStorage.UpdateStatus(ctx, order.ID, domain.StatusPayed); err != nil {
			return fmt.Errorf("orderStorage.Save: %w", err)
		}

		event, err := domain.NewOrderEvent(order.ID, domain.StatusPayed)
		if err != nil {
			return fmt.Errorf("domain.NewOrderEvent: %w", err)
		}

		if err = l.outboxStorage.CreateEvent(ctx, event); err != nil {
			return fmt.Errorf("outboxStorage.CreateEvent: %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("transaction.Wrap: %w", err)
	}

	return nil
}
