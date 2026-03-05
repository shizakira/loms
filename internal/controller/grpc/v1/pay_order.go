package v1

import (
	"context"

	pb "github.com/shizakira/loms/gen/grpc/loms_v1"
	"github.com/shizakira/loms/internal/dto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h Handlers) OrderPay(ctx context.Context, req *pb.OrderPayRequest) (*emptypb.Empty, error) {
	input := dto.PayOrderInput{
		OrderID: int(req.GetOrderId()),
	}

	if err := h.usecase.PayOrder(ctx, input); err != nil {
		return nil, GRPCError(ctx, err)
	}

	return &emptypb.Empty{}, nil
}
