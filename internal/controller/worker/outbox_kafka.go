package worker

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/shizakira/loms/internal/usecase"
)

type OutboxKafkaConfig struct {
	Limit    int `envconfig:"OUTBOX_KAFKA_WORKER_LIMIT" default:"10"`
	Duration int `envconfig:"OUTBOX_KAFKA_WORKER_DURATION" default:"5"`
}

type OutboxKafka struct {
	config  OutboxKafkaConfig
	usecase *usecase.Loms
	stop    chan struct{}
	done    chan struct{}
}

func NewOutboxKafka(uc *usecase.Loms, c OutboxKafkaConfig) *OutboxKafka {
	w := &OutboxKafka{
		config:  c,
		usecase: uc,
		stop:    make(chan struct{}),
		done:    make(chan struct{}),
	}

	go w.run()

	return w
}

func (w *OutboxKafka) run() {
	log.Info().Msg("outbox kafka worker: started")

FOR:
	for {
		count, err := w.usecase.OutboxReadAndProduce(context.Background(), w.config.Limit)
		if err != nil {
			log.Error().Err(err).Msg("outbox kafka worker: read and produce failed")
		}

		log.Info().Int("count", count).Msg("outbox kafka worker: read and produce")

		var duration time.Duration
		if count < w.config.Limit {
			duration = time.Duration(w.config.Duration) * time.Second
			log.Info().Msgf("outbox kafka worker: sleeping %s", duration.String())
		}

		select {
		case <-w.stop:
			break FOR
		case <-time.After(duration):
		}
	}

	close(w.done)
}

func (w *OutboxKafka) Close() {
	log.Info().Msg("outbox kafka worker: closing")

	close(w.stop)

	<-w.done

	log.Info().Msg("outbox kafka worker: closed")
}
