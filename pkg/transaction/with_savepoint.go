package transaction

import (
	"context"
	"fmt"
)

func WithSavepoint(ctx context.Context, savepointName string, fn func(context.Context) error) error {
	if IsUnitTest {
		return fn(ctx)
	}

	tx := TryExtractTX(ctx)

	if _, err := tx.Exec(ctx, "SAVEPOINT "+savepointName); err != nil {
		return fmt.Errorf("tx.Exec savepoint: %w", err)
	}

	if err := fn(ctx); err != nil {
		if _, rbErr := tx.Exec(ctx, "ROLLBACK TO SAVEPOINT "+savepointName); rbErr != nil {
			return fmt.Errorf("tx.Exec rollback to savepoint: %w", rbErr)
		}
		return err
	}

	if _, err := tx.Exec(ctx, "RELEASE SAVEPOINT "+savepointName); err != nil {
		return fmt.Errorf("tx.Exec savepoint: %w", err)
	}

	return nil
}
