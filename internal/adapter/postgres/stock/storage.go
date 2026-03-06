package stock

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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

func (s *Storage) GetBySkuID(ctx context.Context, sku int) (domain.Stock, error) {
	q := s.query(ctx)

	row, err := q.GetStockBySku(ctx, int64(sku))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Stock{}, domain.ErrSkuNotFound
		}
		return domain.Stock{}, fmt.Errorf("q.GetStockBySku: %w", err)
	}

	return domain.Stock{
		SkuID:      int(row.SkuID),
		TotalCount: int(row.TotalCount),
		Reserved:   int(row.Reserved),
	}, nil
}

func (s *Storage) DecreaseReserved(ctx context.Context, skuID int, count int) error {
	q := s.query(ctx)

	err := q.DecreaseReservedStock(ctx, sqlc.DecreaseReservedStockParams{
		SkuID: int64(skuID),
		Count: int64(count),
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23514" {
			return domain.ErrInsufficientStock
		}

		return fmt.Errorf("q.DecreaseReserved: %w", err)

	}

	return nil
}

func (s *Storage) Reserve(ctx context.Context, skuID int, count int) error {
	q := s.query(ctx)

	affected, err := q.IncreaseReservedStock(ctx, sqlc.IncreaseReservedStockParams{
		Count: int64(count),
		SkuID: int64(skuID),
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23514" {
			return domain.ErrInsufficientStock
		}

		return fmt.Errorf("q.IncreaseReservedStock: %w", err)

	}
	if affected == 0 {
		return domain.ErrSkuNotFound
	}

	return nil
}

func (s *Storage) DecreaseReserveAndTotalCount(ctx context.Context, skuID int, count int) error {
	q := s.query(ctx)

	err := q.DecreaseReserveAndTotalCountStock(ctx, sqlc.DecreaseReserveAndTotalCountStockParams{
		Count: int64(count),
		SkuID: int64(skuID),
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23514" {
			return domain.ErrInsufficientStock
		}

		return fmt.Errorf("q.IncreaseReservedStock: %w", err)

	}

	return nil
}
