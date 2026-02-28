package http_gateway

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/flowchartsman/swaggerui"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	pb "github.com/shizakira/loms/gen/grpc/loms_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	Port     string `default:"8080" envconfig:"HTTP_PORT"`
	GRPCAddr string `default:"localhost:50051" envconfig:"GRPC_ADDR"`
}

type Server struct {
	server *http.Server
}

func New(ctx context.Context, c Config) (*Server, error) {
	gwMux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := pb.RegisterLomsHandlerFromEndpoint(ctx, gwMux, c.GRPCAddr, opts); err != nil {
		return nil, fmt.Errorf("RegisterLomsHandlerFromEndpoint: %w", err)
	}

	spec, _ := getSpec()

	mux := http.NewServeMux()
	mux.Handle("/", gwMux)
	mux.Handle("/swagger/", http.StripPrefix("/swagger", swaggerui.Handler(spec)))

	s := &Server{
		server: &http.Server{
			Addr:    ":" + c.Port,
			Handler: mux,
		},
	}

	go s.start()

	log.Info().Msg("http gateway: started on port: " + c.Port)

	return s, nil
}

func (s *Server) start() {
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error().Err(err).Msg("http gateway: ListenAndServe")
	}
}

func (s *Server) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("http gateway: Shutdown")
	}
	log.Info().Msg("http gateway: closed")
}

func getSpec() ([]byte, error) {
	spec, err := os.ReadFile("gen/openapiv2/loms_v1.swagger.json")
	if err != nil {
		return nil, fmt.Errorf("os.ReadFile swagger: %w", err)
	}

	return spec, nil
}
