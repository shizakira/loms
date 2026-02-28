package usecase

import (
	"context"
	"fmt"

	"github.com/shizakira/loms/internal/domain"
	"github.com/shizakira/loms/internal/dto"
)

//go:generate mockery

type OrderStorage interface {
	Create(ctx context.Context, order domain.Order) (domain.Order, error)
	GetByID(ctx context.Context, orderID int) (domain.Order, error)
	Save(ctx context.Context, order domain.Order) error
}
type StockStorage interface {
	GetBySkuIDs(ctx context.Context, skus []uint32) ([]domain.Stock, error)
	GetBySkuID(ctx context.Context, skus uint32) (domain.Stock, error)
	Save(ctx context.Context, stock domain.Stock) error
}

type Loms struct {
	orderStorage OrderStorage
	stockStorage StockStorage
}

func New(orderStorage OrderStorage, stockStorage StockStorage) *Loms {
	return &Loms{orderStorage: orderStorage, stockStorage: stockStorage}
}

func (l *Loms) reserveStocks(
	ctx context.Context,
	items []domain.OrderItem,
	stocks []domain.Stock,
) error {
	stockMap := make(map[uint32]domain.Stock, len(stocks))
	for _, s := range stocks {
		stockMap[s.SkuID] = s
	}

	for _, item := range items {
		stock, ok := stockMap[item.Sku]
		if !ok {
			return fmt.Errorf("id %d: %w", item.Sku, domain.ErrSkuNotFound)
		}
		if stock.TotalCount-stock.Reserved < uint64(item.Count) {
			return domain.ErrInsufficientStock
		}
		stock.Reserved += uint64(item.Count)
		if err := l.stockStorage.Save(ctx, stock); err != nil {
			return fmt.Errorf("stockStorage.Save: %w", err)
		}
	}

	return nil
}

func (l *Loms) aggregateItems(items []dto.OrderItem) []domain.OrderItem {
	m := make(map[uint32]uint16, len(items))

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
