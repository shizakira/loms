package domain

import "errors"

var ErrInsufficientStock = errors.New("insufficient stock for reservation")
var ErrSkuNotFound = errors.New("sku not found")
