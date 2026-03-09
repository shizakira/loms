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

type OutboxStorage interface {
	CreateEvent(ctx context.Context, events ...domain.Event) error
	FetchMsgs(ctx context.Context, limit int) ([]domain.Event, error)
}

type Producer interface {
	EmitEvents(ctx context.Context, events ...domain.Event) error
}

type Loms struct {
	orderStorage  OrderStorage
	stockStorage  StockStorage
	outboxStorage OutboxStorage
	producer      Producer
}

func New(
	orderStorage OrderStorage,
	stockStorage StockStorage,
	outboxStorage OutboxStorage,
	producer Producer,
) *Loms {
	return &Loms{
		orderStorage:  orderStorage,
		stockStorage:  stockStorage,
		outboxStorage: outboxStorage,
		producer:      producer,
	}
}
