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
		order, err := l.orderStorage.GetByID(ctx, input.OrderID)
		if err != nil {
			return fmt.Errorf("orderStorage.GetByID: %w", err)
		}

		if order.Status != domain.StatusAwaitingPayment {
			return domain.ErrInvalidOrderStatus
		}

		skuIDs := make([]uint32, 0, len(order.Items))
		for _, item := range order.Items {
			skuIDs = append(skuIDs, item.Sku)
		}
		stocks, err := l.stockStorage.GetBySkuIDs(ctx, skuIDs)
		if err != nil {
			return fmt.Errorf("stockStorage.GetBySkuIDs: %w", err)
		}

		stockMap := make(map[uint32]domain.Stock, len(stocks))
		for _, s := range stocks {
			stockMap[s.SkuID] = s
		}

		for _, item := range order.Items {
			stock := stockMap[item.Sku]
			stock.TotalCount -= uint64(item.Count)
			stock.Reserved -= uint64(item.Count)
			if err = l.stockStorage.Save(ctx, stock); err != nil {
				return fmt.Errorf("stockStorage.Save: %w", err)
			}
		}

		order.Status = domain.StatusPayed
		if err := l.orderStorage.Save(ctx, order); err != nil {
			return fmt.Errorf("orderStorage.Save: %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("transaction.Wrap: %w", err)
	}

	return nil
}
