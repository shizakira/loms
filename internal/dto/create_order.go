package dto

type OrderItem struct {
	Sku   uint32
	Count uint16
}

type CreateOrderInput struct {
	User  int
	Items []OrderItem
}

type CreateOrderOutput struct {
	OrderID int
}
