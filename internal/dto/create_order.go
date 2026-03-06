package dto

type OrderItem struct {
	Sku   int
	Count int
}

type CreateOrderInput struct {
	User  int
	Items []OrderItem
}

type CreateOrderOutput struct {
	OrderID int
}
