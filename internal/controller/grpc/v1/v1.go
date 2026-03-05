package v1

import (
	"context"
	"errors"

	pb "github.com/shizakira/loms/gen/grpc/loms_v1"
	"github.com/shizakira/loms/internal/domain"
	"github.com/shizakira/loms/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/shizakira/loms/internal/usecase"
)

type Handlers struct {
	pb.UnimplementedLomsServer
	usecase *usecase.Loms
}

func New(uc *usecase.Loms) *Handlers {
	return &Handlers{usecase: uc}
}

func GRPCError(ctx context.Context, err error) error {
	ctxErr, ok := ctx.Value(logger.ContextErrKey{}).(*error)
	if ok {
		*ctxErr = err
	}

	var domainErr *domain.Error
	if errors.As(err, &domainErr) {
		return status.Error(domainErr.GRPCCode, domainErr.Message)
	}

	return status.Error(codes.Internal, "internal error")
}
