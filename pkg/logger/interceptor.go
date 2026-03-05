package logger

import (
	"context"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type ContextErrKey struct{}

func Interceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	var ctxErr error
	ctx = context.WithValue(ctx, ContextErrKey{}, &ctxErr)

	event := log.Info()

	resp, err := handler(ctx, req)
	if ctxErr != nil {
		event = log.Error().Err(ctxErr)
	} else if err != nil {
		event = log.Error().Err(err)
	}

	event.
		Str("code", status.Code(ctxErr).String()).
		Str("grpc_method", info.FullMethod).
		Send()

	return resp, err
}
