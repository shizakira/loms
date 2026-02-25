package domain

type OrderStatus string

const (
	StatusNew             OrderStatus = "new"
	StatusAwaitingPayment OrderStatus = "awaiting payment"
	StatusFailed          OrderStatus = "failed"
	StatusPayed           OrderStatus = "payed"
	StatusCancelled       OrderStatus = "cancelled"
)

type OrderItem struct {
	SkuID uint32
	Count uint16
}

type Order struct {
	ID     int
	Status OrderStatus
	UserID int
	Items  []OrderItem
}
