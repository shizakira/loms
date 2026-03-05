package dto

import "github.com/shizakira/loms/internal/domain"

type GetOrderInfoInput struct {
	OrderID int
}

type GetOrderInfoOutput struct {
	domain.Order
}
