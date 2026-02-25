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
	skuIDs := make([]uint32, 0, len(input.Items))
	for _, item := range items {
		skuIDs = append(skuIDs, item.SkuID)
	}

	order := domain.Order{
		Status: domain.StatusNew,
		UserID: input.UserID,
		Items:  items,
	}

	order, err := l.orderStorage.Create(ctx, order)
	if err != nil {
		return dto.CreateOrderOutput{}, fmt.Errorf("orderStorage.Create: %w", err)
	}

	err = transaction.Wrap(ctx, func(ctx context.Context) error {
		stocks, err := l.stockStorage.GetBySkuIDs(ctx, skuIDs)
		if err != nil {
			return fmt.Errorf("stockStorage.GetBySkuIDs: %w", err)
		}

		if err = l.reserveStocks(ctx, items, stocks); err != nil {
			order.Status = domain.StatusFailed
			if err := l.orderStorage.Save(ctx, order); err != nil {
				return fmt.Errorf("orderStorage.Save: %w", err)
			}

			return fmt.Errorf("usecase.reserveStocks: %w", err)
		}

		order.Status = domain.StatusAwaitingPayment
		if err := l.orderStorage.Save(ctx, order); err != nil {
			return fmt.Errorf("orderStorage.Save: %w", err)
		}

		return nil
	})

	if err != nil {
		return dto.CreateOrderOutput{}, fmt.Errorf("transaction.Wrap: %w", err)
	}

	return dto.CreateOrderOutput{OrderID: order.ID}, nil
}
