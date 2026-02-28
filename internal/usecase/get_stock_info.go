package usecase

import (
	"context"
	"fmt"

	"github.com/shizakira/loms/internal/dto"
)

func (l *Loms) GetStockInfo(ctx context.Context, input dto.GetStockInfoInput) (dto.GetStockInfoOutput, error) {
	stock, err := l.stockStorage.GetBySkuID(ctx, input.Sku)
	if err != nil {
		return dto.GetStockInfoOutput{}, fmt.Errorf("stockStorage.GetBySkuID: %w", err)
	}
	count := stock.TotalCount - stock.Reserved

	return dto.GetStockInfoOutput{Count: int(count)}, nil
}
