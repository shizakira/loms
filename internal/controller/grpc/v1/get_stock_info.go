package v1

import (
	"context"

	pb "github.com/shizakira/loms/gen/grpc/loms_v1"
	"github.com/shizakira/loms/internal/dto"
)

func (h Handlers) StocksInfo(ctx context.Context, req *pb.StocksInfoRequest) (*pb.StocksInfoResponse, error) {
	input := dto.GetStockInfoInput{
		Sku: int(req.GetSku()),
	}

	stock, err := h.usecase.GetStockInfo(ctx, input)
	if err != nil {
		return nil, GRPCError(ctx, err)
	}

	return &pb.StocksInfoResponse{
		Count: uint64(stock.Count),
	}, nil
}
