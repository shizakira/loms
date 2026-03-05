//go:build integration

package test

import (
	"context"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	pb "github.com/shizakira/loms/gen/grpc/loms_v1"
	"github.com/shizakira/loms/internal/app"
	"github.com/shizakira/loms/internal/config"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Test_Integration(t *testing.T) {
	suite.Run(t, &Suite{})
}

type Suite struct {
	suite.Suite
	*require.Assertions
	pb.LomsClient
}

func (s *Suite) SetupSuite() {
	s.Assertions = s.Require()

	godotenv.Overload(".env.test")
	c, err := config.NewTest()
	s.NoError(err)

	s.ResetMigrations()

	log.Logger = zerolog.Nop()

	go func() {
		err = app.Run(context.Background(), c)
		s.NoError(err)
	}()

	time.Sleep(200 * time.Millisecond)

	conn, err := grpc.NewClient(
		c.GRPC.Addr+":"+c.GRPC.Port,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	s.NoError(err)
	s.LomsClient = pb.NewLomsClient(conn)
}
