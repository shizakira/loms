package usecase

import (
	"context"
	"fmt"

	"github.com/shizakira/loms/internal/dto"
)

func (l *Loms) GetOrderInfo(ctx context.Context, input dto.GetOrderInfoInput) (dto.GetOrderInfoOutput, error) {
	order, err := l.orderStorage.GetByID(ctx, input.OrderID, false)
	if err != nil {
		return dto.GetOrderInfoOutput{}, fmt.Errorf("orderStorage.GetByID: %w", err)
	}

	return dto.GetOrderInfoOutput{Order: order}, nil
}
