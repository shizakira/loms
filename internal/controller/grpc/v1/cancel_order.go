package v1

import (
	"context"

	pb "github.com/shizakira/loms/gen/grpc/loms_v1"
	"github.com/shizakira/loms/internal/dto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h Handlers) OrderCancel(ctx context.Context, req *pb.OrderCancelRequest) (*emptypb.Empty, error) {
	input := dto.CancelOrderInput{
		OrderID: int(req.GetOrderId()),
	}

	if err := h.usecase.CancelOrder(ctx, input); err != nil {
		return nil, GRPCError(ctx, err)
	}

	return &emptypb.Empty{}, nil
}
