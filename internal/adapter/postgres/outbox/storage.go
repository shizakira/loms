package outbox

import (
	"context"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/shizakira/loms/internal/domain"
	"github.com/shizakira/loms/pkg/transaction"
)

type Storage struct{}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) CreateEvent(ctx context.Context, events ...domain.Event) error {
	if len(events) == 0 {
		return nil
	}

	batch := make([]any, 0, len(events))
	for _, event := range events {
		batch = append(batch, goqu.Record{
			"topic": event.Topic,
			"key":   event.Key,
			"value": event.Value,
		})
	}

	sql, _, err := goqu.Insert("outbox").Rows(batch...).ToSQL()
	if err != nil {
		return fmt.Errorf("goqu.Insert.ToSQL: %w", err)
	}

	txOrPool := transaction.TryExtractTX(ctx)

	if _, err = txOrPool.Exec(ctx, sql); err != nil {
		return fmt.Errorf("txOrPool.Exec: %w", err)
	}

	return nil
}

func (s *Storage) FetchMsgs(ctx context.Context, limit int) ([]domain.Event, error) {
	const sql = `WITH taken as (SELECT id, topic, key, value
				 	FROM outbox
					ORDER BY created_at
					LIMIT $1 FOR UPDATE SKIP LOCKED)
				 DELETE
				 FROM outbox
				 WHERE id in (SELECT id FROM taken)
				 RETURNING topic, key, value;`

	txOrPool := transaction.TryExtractTX(ctx)
	rows, err := txOrPool.Query(ctx, sql, limit)
	if err != nil {
		return nil, fmt.Errorf("txOrPool.Query: %w", err)
	}
	defer rows.Close()

	events := make([]domain.Event, 0, limit)
	for rows.Next() {
		var event domain.Event

		if err = rows.Scan(&event.Topic, &event.Key, &event.Value); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		events = append(events, event)
	}

	return events, nil
}
