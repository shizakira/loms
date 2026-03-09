package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/shizakira/loms/internal/adapter/kafka_producer"
	"github.com/shizakira/loms/internal/controller/http_gateway"
	"github.com/shizakira/loms/internal/controller/worker"
	"github.com/shizakira/loms/pkg/logger"

	"github.com/shizakira/loms/internal/controller/grpc"
	"github.com/shizakira/loms/pkg/postgres"
)

type App struct {
	Name    string `envconfig:"APP_NAME"    required:"true"`
	Version string `envconfig:"APP_VERSION" required:"true"`
}
type Config struct {
	App           App
	Postgres      postgres.Config
	Logger        logger.Config
	GRPC          grpc.Config
	HTTP          http_gateway.Config
	KafkaProducer kafka_producer.Config
	OutboxKafka   worker.OutboxKafkaConfig
}

func New() (Config, error) {
	return load(".env")
}

func NewTest() (Config, error) {
	return load(".env.test")
}

func load(env string) (Config, error) {
	var config Config

	if err := godotenv.Load(env); err != nil {
		return config, fmt.Errorf("godotenv.Load: %w", err)
	}

	if err := envconfig.Process("", &config); err != nil {
		return config, fmt.Errorf("envconfig.Process: %w", err)
	}

	return config, nil
}
