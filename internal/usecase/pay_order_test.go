package usecase_test

import (
	"context"
	"testing"

	"github.com/shizakira/loms/internal/domain"
	"github.com/shizakira/loms/internal/dto"
	"github.com/shizakira/loms/internal/usecase"
	"github.com/shizakira/loms/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestLoms_PayOrder_Success(t *testing.T) {
	ctx := context.Background()
	input := dto.PayOrderInput{OrderID: 1}

	order := domain.Order{
		ID:     42,
		Status: domain.StatusAwaitingPayment,
		Items:  []domain.OrderItem{{Sku: 1, Count: 2}},
	}

	orderStorage := &mocks.OrderStorage{}
	orderStorage.On("GetByID", ctx, input.OrderID).Return(order, nil).Once()
	orderStorage.On("Save", ctx, mock.Anything).Return(nil).Once()

	stockStorage := &mocks.StockStorage{}
	stockStorage.On("GetBySkuIDs", ctx, mock.Anything).Return([]domain.Stock{
		{SkuID: 1, TotalCount: 10, Reserved: 2},
	}, nil).Once()
	stockStorage.On("Save", ctx, mock.Anything).Return(nil).Once()

	uc := usecase.New(orderStorage, stockStorage)
	err := uc.PayOrder(ctx, input)

	require.NoError(t, err)
	orderStorage.AssertExpectations(t)
	stockStorage.AssertExpectations(t)
}

func TestLoms_PayOrder_ReturnsErrorWhenOrderNotFound(t *testing.T) {
	ctx := context.Background()
	input := dto.PayOrderInput{OrderID: 42}

	orderStorage := &mocks.OrderStorage{}
	orderStorage.On("GetByID", ctx, input.OrderID).Return(domain.Order{}, assert.AnError).Once()

	stockStorage := &mocks.StockStorage{}

	uc := usecase.New(orderStorage, stockStorage)
	err := uc.PayOrder(ctx, input)

	require.ErrorIs(t, err, assert.AnError)
	orderStorage.AssertExpectations(t)
	stockStorage.AssertExpectations(t)
}

func TestLoms_PayOrder_ReturnsErrorWhenGetStocksFails(t *testing.T) {
	ctx := context.Background()
	input := dto.PayOrderInput{OrderID: 42}

	order := domain.Order{
		ID:     42,
		Status: domain.StatusAwaitingPayment,
		Items:  []domain.OrderItem{{Sku: 1, Count: 2}},
	}

	orderStorage := &mocks.OrderStorage{}
	orderStorage.On("GetByID", ctx, input.OrderID).Return(order, nil).Once()

	stockStorage := &mocks.StockStorage{}
	stockStorage.On("GetBySkuIDs", ctx, mock.Anything).Return(nil, assert.AnError).Once()

	uc := usecase.New(orderStorage, stockStorage)
	err := uc.PayOrder(ctx, input)

	require.ErrorIs(t, err, assert.AnError)
	orderStorage.AssertExpectations(t)
	stockStorage.AssertExpectations(t)
}
