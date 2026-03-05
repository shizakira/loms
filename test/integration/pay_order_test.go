//go:build integration

package test

import (
	"context"

	pb "github.com/shizakira/loms/gen/grpc/loms_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Suite) Test_OrderPay_Success() {
	created, err := s.OrderCreate(context.Background(), &pb.OrderCreateRequest{
		User: 1,
		Items: []*pb.Item{
			{Sku: 1076963, Count: 1},
		},
	})
	s.NoError(err)

	_, err = s.OrderPay(context.Background(), &pb.OrderPayRequest{
		OrderId: created.GetOrderId(),
	})
	s.NoError(err)

	info, err := s.OrderInfo(context.Background(), &pb.OrderInfoRequest{
		OrderId: created.GetOrderId(),
	})
	s.NoError(err)
	s.Equal("payed", info.GetStatus())
}

func (s *Suite) Test_OrderPay_NotFound() {
	_, err := s.OrderPay(context.Background(), &pb.OrderPayRequest{
		OrderId: 99999999,
	})
	s.Equal(codes.NotFound, status.Code(err))
}

func (s *Suite) Test_OrderPay_InvalidStatus() {
	created, err := s.OrderCreate(context.Background(), &pb.OrderCreateRequest{
		User: 1,
		Items: []*pb.Item{
			{Sku: 1076963, Count: 1},
		},
	})
	s.NoError(err)

	_, err = s.OrderPay(context.Background(), &pb.OrderPayRequest{OrderId: created.GetOrderId()})
	s.NoError(err)

	_, err = s.OrderPay(context.Background(), &pb.OrderPayRequest{OrderId: created.GetOrderId()})
	s.Equal(codes.FailedPrecondition, status.Code(err))
}
