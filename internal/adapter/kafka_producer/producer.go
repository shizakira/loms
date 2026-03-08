package kafka_producer

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	"github.com/shizakira/loms/internal/domain"
	"github.com/shizakira/loms/pkg/logger"
)

type Config struct {
	Addr []string `envconfig:"KAFKA_WRITER_ADDR"  required:"true"`
}

type Producer struct {
	config Config
	writer *kafka.Writer
}

func New(c Config) *Producer {
	w := &kafka.Writer{
		Addr:         kafka.TCP(c.Addr...),
		Balancer:     &kafka.Hash{},
		RequiredAcks: kafka.RequireAll,
		ErrorLogger:  logger.ErrorLogger(),
		Async:        false,
	}

	return &Producer{
		config: c,
		writer: w,
	}
}

func (p *Producer) EmitEvents(ctx context.Context, events ...domain.Event) error {
	var msgs []kafka.Message

	for _, e := range events {
		msg := kafka.Message{
			Topic: e.Topic,
			Key:   e.Key,
			Value: e.Value,
		}
		msgs = append(msgs, msg)
	}

	if err := p.writer.WriteMessages(ctx, msgs...); err != nil {
		return fmt.Errorf("p.writer.WriteMessages: %w", err)
	}

	return nil
}

func (p *Producer) Close() {
	if err := p.writer.Close(); err != nil {
		log.Error().Err(err).Msg("kafka producer: p.writer.Close")
	}

	log.Info().Msg("kafka producer: closed")
}
