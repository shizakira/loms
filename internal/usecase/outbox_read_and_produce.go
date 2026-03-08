package usecase

import (
	"context"
	"fmt"

	"github.com/shizakira/loms/pkg/transaction"
)

func (l *Loms) OutboxReadAndProduce(ctx context.Context, limit int) (count int, err error) {
	err = transaction.Wrap(ctx, func(ctx context.Context) error {
		msgs, err := l.outboxStorage.FetchMsgs(ctx, limit)
		if err != nil {
			return fmt.Errorf("outboxStorage.FetchMsgs: %w", err)
		}

		count = len(msgs)
		if err = l.producer.EmitEvents(ctx, msgs...); err != nil {
			return fmt.Errorf("producer.Produce: %w", err)
		}

		return nil
	})
	if err != nil {
		return count, fmt.Errorf("transaction.Wrap: %w", err)
	}

	return count, nil
}
