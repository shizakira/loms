package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
	"github.com/shizakira/loms/internal/adapter/kafka_producer"
	order_storage "github.com/shizakira/loms/internal/adapter/postgres/order"
	outbox_storage "github.com/shizakira/loms/internal/adapter/postgres/outbox"
	stock_storage "github.com/shizakira/loms/internal/adapter/postgres/stock"
	"github.com/shizakira/loms/internal/config"
	"github.com/shizakira/loms/internal/controller/grpc"
	"github.com/shizakira/loms/internal/controller/http_gateway"
	"github.com/shizakira/loms/internal/controller/worker"
	"github.com/shizakira/loms/internal/usecase"
	pgpool "github.com/shizakira/loms/pkg/postgres"
	"github.com/shizakira/loms/pkg/transaction"
)

func Run(ctx context.Context, c config.Config) error {
	pgPool, err := pgpool.New(ctx, c.Postgres)
	if err != nil {
		return fmt.Errorf("postgres.New: %w", err)
	}

	transaction.Init(pgPool)

	kafkaProducer := kafka_producer.New(c.KafkaProducer)

	uc := usecase.New(
		order_storage.New(),
		stock_storage.New(),
		outbox_storage.New(),
		kafkaProducer,
	)

	outboxKafkaWorker := worker.NewOutboxKafka(uc, c.OutboxKafka)

	grpcServer, err := grpc.New(c.GRPC, uc)
	if err != nil {
		return fmt.Errorf("grpc.New: %w", err)
	}

	httpGatewayServer, err := http_gateway.New(ctx, c.HTTP)
	if err != nil {
		return fmt.Errorf("http_gateway.New: %w", err)
	}

	log.Info().Msg("app started!")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig

	log.Info().Msg("app got signal to stop")

	httpGatewayServer.Close()
	grpcServer.Close()
	outboxKafkaWorker.Close()

	kafkaProducer.Close()
	pgPool.Close()

	log.Info().Msg("app stopped")

	return nil

}
