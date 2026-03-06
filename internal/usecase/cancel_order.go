package usecase

import (
	"context"
	"fmt"

	"github.com/shizakira/loms/internal/domain"
	"github.com/shizakira/loms/internal/dto"
	"github.com/shizakira/loms/pkg/transaction"
)

func (l *Loms) CancelOrder(ctx context.Context, input dto.CancelOrderInput) error {
	err := transaction.Wrap(ctx, func(ctx context.Context) error {
		order, err := l.orderStorage.GetByID(ctx, input.OrderID, true)
		if err != nil {
			return fmt.Errorf("orderStorage.GetByID: %w", err)
		}

		if order.Status != domain.StatusNew && order.Status != domain.StatusAwaitingPayment {
			return domain.ErrInvalidOrderStatus
		}

		for _, item := range order.Items {
			if err := l.stockStorage.DecreaseReserved(ctx, item.Sku, item.Count); err != nil {
				return fmt.Errorf("stockStorage.DecreaseReserved: %w", err)
			}
		}

		if err := l.orderStorage.UpdateStatus(ctx, input.OrderID, domain.StatusCancelled); err != nil {
			return fmt.Errorf("orderStorage.Save: %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("transaction.Wrap: %w", err)
	}

	return nil
}
