package main

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/shizakira/loms/internal/app"
	"github.com/shizakira/loms/internal/config"
	"github.com/shizakira/loms/pkg/logger"
)

func main() {
	c, err := config.New()
	if err != nil {
		log.Fatal().Err(err).Msg("config.New")
	}

	logger.Init(c.Logger)

	ctx := context.Background()

	if err = app.Run(ctx, c); err != nil {
		log.Error().Err(err).Msg("app.Run")
	}
}
