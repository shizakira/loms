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
	Sku   int `json:"sku"`
	Count int `json:"count"`
}

type Order struct {
	ID     int
	Status OrderStatus
	User   int
	Items  []OrderItem
}
