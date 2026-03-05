//go:build integration

package test

import (
	"context"

	pb "github.com/shizakira/loms/gen/grpc/loms_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Suite) Test_StocksInfo_Success() {
	resp, err := s.StocksInfo(context.Background(), &pb.StocksInfoRequest{
		Sku: 1076963,
	})
	s.NoError(err)
	s.Greater(resp.GetCount(), uint64(0))
}

func (s *Suite) Test_StocksInfo_NotFound() {
	_, err := s.StocksInfo(context.Background(), &pb.StocksInfoRequest{
		Sku: 99999999,
	})
	s.Equal(codes.NotFound, status.Code(err))
}

func (s *Suite) Test_StocksInfo_InvalidSku() {
	_, err := s.StocksInfo(context.Background(), &pb.StocksInfoRequest{
		Sku: 0,
	})
	s.Equal(codes.InvalidArgument, status.Code(err))
}
