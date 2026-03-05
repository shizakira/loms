package usecase_test

import (
	"context"
	"os"
	"testing"

	"github.com/shizakira/loms/internal/domain"
	"github.com/shizakira/loms/internal/dto"
	"github.com/shizakira/loms/internal/usecase"
	"github.com/shizakira/loms/internal/usecase/mocks"
	"github.com/shizakira/loms/pkg/transaction"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	transaction.IsUnitTest = true
	os.Exit(m.Run())
}

func TestLoms_CreateOrder_Success(t *testing.T) {
	ctx := context.Background()
	input := dto.CreateOrderInput{
		User: 1,
		Items: []dto.OrderItem{
			{Sku: 1, Count: 1},
			{Sku: 1, Count: 1},
		},
	}

	createdOrder := domain.Order{ID: 42, Status: domain.StatusNew, User: input.User}

	orderStorage := &mocks.OrderStorage{}
	orderStorage.On("Create", ctx, mock.Anything).Return(createdOrder, nil).Once()
	orderStorage.On("Save", ctx, mock.Anything).Return(nil).Once()

	stockStorage := &mocks.StockStorage{}
	stockStorage.On("GetBySkuIDs", ctx, mock.Anything).
		Return([]domain.Stock{{SkuID: 1, TotalCount: 10}}, nil).Once()
	stockStorage.On("Save", ctx, mock.Anything).Return(nil).Once()

	uc := usecase.New(orderStorage, stockStorage)
	out, err := uc.CreateOrder(ctx, input)

	require.NoError(t, err)
	require.Equal(t, createdOrder.ID, out.OrderID)
	orderStorage.AssertExpectations(t)
	stockStorage.AssertExpectations(t)
}

func TestLoms_CreateOrder_ReturnsErrorWhenCreateFails(t *testing.T) {
	ctx := context.Background()
	input := dto.CreateOrderInput{
		User:  1,
		Items: []dto.OrderItem{{Sku: 1, Count: 1}},
	}

	orderStorage := &mocks.OrderStorage{}
	orderStorage.On("Create", ctx, mock.Anything).Return(domain.Order{}, assert.AnError).Once()

	stockStorage := &mocks.StockStorage{}

	uc := usecase.New(orderStorage, stockStorage)
	_, err := uc.CreateOrder(ctx, input)

	require.ErrorIs(t, err, assert.AnError)
	orderStorage.AssertExpectations(t)
	stockStorage.AssertExpectations(t)
}

func TestLoms_CreateOrder_ReturnsErrorWhenInsufficientStock(t *testing.T) {
	ctx := context.Background()
	input := dto.CreateOrderInput{
		User:  1,
		Items: []dto.OrderItem{{Sku: 1, Count: 10}},
	}

	createdOrder := domain.Order{ID: 42, Status: domain.StatusNew, User: input.User}

	orderStorage := &mocks.OrderStorage{}
	orderStorage.On("Create", ctx, mock.Anything).Return(createdOrder, nil).Once()
	orderStorage.On("Save", ctx, mock.Anything).Return(nil).Once()

	stockStorage := &mocks.StockStorage{}
	stockStorage.On("GetBySkuIDs", ctx, mock.Anything).
		Return([]domain.Stock{{SkuID: 1, TotalCount: 5}}, nil).Once()

	uc := usecase.New(orderStorage, stockStorage)
	_, err := uc.CreateOrder(ctx, input)

	require.ErrorIs(t, err, domain.ErrInsufficientStock)
	orderStorage.AssertExpectations(t)
	stockStorage.AssertExpectations(t)
}
