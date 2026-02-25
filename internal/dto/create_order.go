package dto

type OrderItem struct {
	SkuID uint32
	Count uint16
}

type CreateOrderInput struct {
	UserID int
	Items  []OrderItem
}

type CreateOrderOutput struct {
	OrderID int
}
