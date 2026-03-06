package v1

import (
	"context"

	pb "github.com/shizakira/loms/gen/grpc/loms_v1"
	"github.com/shizakira/loms/internal/dto"
)

func (h Handlers) OrderInfo(ctx context.Context, req *pb.OrderInfoRequest) (*pb.OrderInfoResponse, error) {
	input := dto.GetOrderInfoInput{
		OrderID: int(req.GetOrderId()),
	}

	order, err := h.usecase.GetOrderInfo(ctx, input)
	if err != nil {
		return nil, GRPCError(ctx, err)
	}

	items := make([]*pb.Item, 0, len(order.Items))
	for _, item := range order.Items {
		items = append(items, &pb.Item{
			Sku:   uint32(item.Sku),
			Count: uint32(item.Count),
		})
	}
	return &pb.OrderInfoResponse{
		Status: string(order.Status),
		User:   int64(order.User),
		Items:  items,
	}, nil
}
