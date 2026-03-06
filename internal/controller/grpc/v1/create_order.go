package v1

import (
	"context"

	pb "github.com/shizakira/loms/gen/grpc/loms_v1"
	"github.com/shizakira/loms/internal/dto"
)

func (h Handlers) OrderCreate(ctx context.Context, req *pb.OrderCreateRequest) (*pb.OrderCreateResponse, error) {
	items := make([]dto.OrderItem, 0, len(req.GetItems()))
	for _, item := range req.GetItems() {
		items = append(items, dto.OrderItem{
			Sku:   int(item.GetSku()),
			Count: int(item.GetCount()),
		})
	}
	input := dto.CreateOrderInput{
		User:  int(req.GetUser()),
		Items: items,
	}
	output, err := h.usecase.CreateOrder(ctx, input)
	if err != nil {
		return nil, GRPCError(ctx, err)
	}

	return &pb.OrderCreateResponse{OrderId: int64(output.OrderID)}, nil

}
