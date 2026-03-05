package grpc

import (
	"fmt"
	"net"

	"buf.build/go/protovalidate"
	protovalidate_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"github.com/rs/zerolog/log"
	pb "github.com/shizakira/loms/gen/grpc/loms_v1"
	loms_handlers_v1 "github.com/shizakira/loms/internal/controller/grpc/v1"
	"github.com/shizakira/loms/internal/usecase"
	"github.com/shizakira/loms/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type Config struct {
	Addr string `default:"localhost" envconfig:"GRPC_ADDR"`
	Port string `default:"50051" envconfig:"GRPC_PORT"`
}

type Server struct {
	server *grpc.Server
}

func New(c Config, uc *usecase.Loms) (*Server, error) {
	validator, err := protovalidate.New()
	if err != nil {
		return nil, fmt.Errorf("protovalidate.New: %w", err)
	}

	s := grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.ChainUnaryInterceptor(
			protovalidate_middleware.UnaryServerInterceptor(validator),
			logger.Interceptor,
		),
	)

	reflection.Register(s)

	v1 := loms_handlers_v1.New(uc)
	pb.RegisterLomsServer(s, v1)

	if err := start(s, c.Addr, c.Port); err != nil {
		return nil, fmt.Errorf("start: %w", err)
	}

	return &Server{server: s}, nil

}

func start(server *grpc.Server, addr string, port string) error {
	conn, err := net.Listen("tcp", net.JoinHostPort(addr, port))
	if err != nil {
		return fmt.Errorf("net.Listen: %w", err)
	}

	go func() {
		if err = server.Serve(conn); err != nil {
			log.Error().Err(err).Msg("grpc server: Serve")
		}
	}()

	log.Info().Msg("grpc server: started on port: " + port)

	return nil
}

func (s *Server) Close() {
	s.server.GracefulStop()

	log.Info().Msg("grpc server: closed")
}
