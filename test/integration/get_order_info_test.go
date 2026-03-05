//go:build integration

package test

import (
	"context"

	pb "github.com/shizakira/loms/gen/grpc/loms_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Suite) Test_OrderInfo_Success() {
	created, err := s.OrderCreate(context.Background(), &pb.OrderCreateRequest{
		User: 1,
		Items: []*pb.Item{
			{Sku: 1076963, Count: 1},
		},
	})
	s.NoError(err)

	resp, err := s.OrderInfo(context.Background(), &pb.OrderInfoRequest{
		OrderId: created.GetOrderId(),
	})
	s.NoError(err)
	s.Equal("awaiting_payment", resp.GetStatus())
	s.Equal(int64(1), resp.GetUser())
}

func (s *Suite) Test_OrderInfo_NotFound() {
	_, err := s.OrderInfo(context.Background(), &pb.OrderInfoRequest{
		OrderId: 99999999,
	})
	s.Equal(codes.NotFound, status.Code(err))
}

func (s *Suite) Test_OrderInfo_InvalidOrderId() {
	_, err := s.OrderInfo(context.Background(), &pb.OrderInfoRequest{
		OrderId: 0,
	})
	s.Equal(codes.InvalidArgument, status.Code(err))
}
