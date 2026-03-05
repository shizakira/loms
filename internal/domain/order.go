package domain

type OrderStatus string

const (
	StatusNew             OrderStatus = "new"
	StatusAwaitingPayment OrderStatus = "awaiting_payment"
	StatusFailed          OrderStatus = "failed"
	StatusPayed           OrderStatus = "payed"
	StatusCancelled       OrderStatus = "cancelled"
)

type OrderItem struct {
	Sku   uint32 `json:"sku"`
	Count uint16 `json:"count"`
}

type Order struct {
	ID     int
	Status OrderStatus
	User   int
	Items  []OrderItem
}
