package domain

import "google.golang.org/grpc/codes"

type Error struct {
	GRPCCode codes.Code
	Message  string
}

func (e *Error) Error() string {
	return e.Message
}

var (
	ErrSkuNotFound        = &Error{GRPCCode: codes.NotFound, Message: "sku not found"}
	ErrOrderNotFound      = &Error{GRPCCode: codes.NotFound, Message: "order not found"}
	ErrInsufficientStock  = &Error{GRPCCode: codes.FailedPrecondition, Message: "insufficient stock"}
	ErrInvalidOrderStatus = &Error{GRPCCode: codes.FailedPrecondition, Message: "invalid order status"}
)
