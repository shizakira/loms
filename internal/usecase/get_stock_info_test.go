package usecase_test

import (
	"context"
	"testing"

	"github.com/shizakira/loms/internal/domain"
	"github.com/shizakira/loms/internal/dto"
	"github.com/shizakira/loms/internal/usecase"
	"github.com/shizakira/loms/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoms_GetStockInfo_Success(t *testing.T) {
	ctx := context.Background()
	input := dto.GetStockInfoInput{Sku: 1}

	stockStorage := &mocks.StockStorage{}
	stockStorage.On("GetBySkuID", ctx, input.Sku).Return(domain.Stock{
		SkuID:      1,
		TotalCount: 10,
		Reserved:   3,
	}, nil).Once()

	uc := usecase.New(&mocks.OrderStorage{}, stockStorage)
	out, err := uc.GetStockInfo(ctx, input)

	require.NoError(t, err)
	require.Equal(t, 7, out.Count)
	stockStorage.AssertExpectations(t)
}

func TestLoms_GetStockInfo_ReturnsErrorWhenSkuNotFound(t *testing.T) {
	ctx := context.Background()
	input := dto.GetStockInfoInput{Sku: 999}

	stockStorage := &mocks.StockStorage{}
	stockStorage.On("GetBySkuID", ctx, input.Sku).Return(domain.Stock{}, assert.AnError).Once()

	uc := usecase.New(&mocks.OrderStorage{}, stockStorage)
	_, err := uc.GetStockInfo(ctx, input)

	require.ErrorIs(t, err, assert.AnError)
	stockStorage.AssertExpectations(t)
}
