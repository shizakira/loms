package usecase

import (
	"context"

	"github.com/shizakira/loms/internal/domain"
)

//go:generate mockery

type OrderStorage interface {
	Create(ctx context.Context, order domain.Order) (int, error)
	GetByID(ctx context.Context, orderID int, pessimistic bool) (domain.Order, error)
	UpdateStatus(ctx context.Context, orderID int, status domain.OrderStatus) error
}
type StockStorage interface {
	GetBySkuID(ctx context.Context, sku int) (domain.Stock, error)
	DecreaseReserved(ctx context.Context, skuID int, count int) error
	DecreaseReserveAndTotalCount(ctx context.Context, skuID int, count int) error
	Reserve(ctx context.Context, skuID int, count int) error
}

type Loms struct {
	orderStorage OrderStorage
	stockStorage StockStorage
}

func New(orderStorage OrderStorage, stockStorage StockStorage) *Loms {
	return &Loms{orderStorage: orderStorage, stockStorage: stockStorage}
}
