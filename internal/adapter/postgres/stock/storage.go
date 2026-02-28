package stock

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/shizakira/loms/internal/adapter/postgres/stock/sqlc"
	"github.com/shizakira/loms/internal/domain"
	"github.com/shizakira/loms/pkg/transaction"
)

type Storage struct{}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) query(ctx context.Context) *sqlc.Queries {
	db := transaction.TryExtractTX(ctx)

	return sqlc.New(db)
}

func (s *Storage) GetBySkuID(ctx context.Context, sku uint32) (domain.Stock, error) {
	q := s.query(ctx)

	row, err := q.GetStockBySku(ctx, int64(sku))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Stock{}, domain.ErrSkuNotFound
		}
		return domain.Stock{}, fmt.Errorf("q.GetStockBySku: %w", err)
	}

	return domain.Stock{
		SkuID:      uint32(row.SkuID),
		TotalCount: uint64(row.TotalCount),
		Reserved:   uint64(row.Reserved),
	}, nil
}

func (s *Storage) GetBySkuIDs(ctx context.Context, skus []uint32) ([]domain.Stock, error) {
	q := s.query(ctx)

	skuIDs := make([]int64, 0, len(skus))
	for _, sku := range skus {
		skuIDs = append(skuIDs, int64(sku))
	}

	rows, err := q.GetStocksBySkus(ctx, skuIDs)
	if err != nil {
		return nil, fmt.Errorf("q.GetStocksBySkus: %w", err)
	}

	stocks := make([]domain.Stock, 0, len(rows))
	for _, row := range rows {
		stocks = append(stocks, domain.Stock{
			SkuID:      uint32(row.SkuID),
			TotalCount: uint64(row.TotalCount),
			Reserved:   uint64(row.Reserved),
		})
	}

	return stocks, nil
}

func (s *Storage) Save(ctx context.Context, stock domain.Stock) error {
	q := s.query(ctx)

	err := q.SaveStock(ctx, sqlc.SaveStockParams{
		SkuID:      int64(stock.SkuID),
		TotalCount: int64(stock.TotalCount),
		Reserved:   int64(stock.Reserved),
	})
	if err != nil {
		return fmt.Errorf("q.SaveStock: %w", err)
	}

	return nil
}
