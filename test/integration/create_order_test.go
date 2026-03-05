//go:build integration

package test

import (
	"context"

	pb "github.com/shizakira/loms/gen/grpc/loms_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Suite) Test_OrderCreate_Success() {
	resp, err := s.OrderCreate(context.Background(), &pb.OrderCreateRequest{
		User: 1,
		Items: []*pb.Item{
			{Sku: 1076963, Count: 1},
		},
	})
	s.NoError(err)
	s.Greater(resp.GetOrderId(), int64(0))
}

func (s *Suite) Test_OrderCreate_InvalidUser() {
	_, err := s.OrderCreate(context.Background(), &pb.OrderCreateRequest{
		User: 0,
		Items: []*pb.Item{
			{Sku: 1076963, Count: 1},
		},
	})
	s.Equal(codes.InvalidArgument, status.Code(err))
}

func (s *Suite) Test_OrderCreate_EmptyItems() {
	_, err := s.OrderCreate(context.Background(), &pb.OrderCreateRequest{
		User:  1,
		Items: []*pb.Item{},
	})
	s.Equal(codes.InvalidArgument, status.Code(err))
}

func (s *Suite) Test_OrderCreate_InsufficientStock() {
	_, err := s.OrderCreate(context.Background(), &pb.OrderCreateRequest{
		User: 1,
		Items: []*pb.Item{
			{Sku: 1076963, Count: 99999},
		},
	})
	s.Equal(codes.FailedPrecondition, status.Code(err))
}

func (s *Suite) Test_OrderCreate_SkuNotFound() {
	_, err := s.OrderCreate(context.Background(), &pb.OrderCreateRequest{
		User: 1,
		Items: []*pb.Item{
			{Sku: 99999999, Count: 1},
		},
	})
	s.Equal(codes.NotFound, status.Code(err))
}
