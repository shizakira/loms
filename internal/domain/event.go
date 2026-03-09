package domain

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type Event struct {
	Topic string
	Key   []byte
	Value []byte
}

type OrderEventOption func(*OrderEvent)

type OrderEvent struct {
	OrderID   int               `json:"order_id"`
	Status    OrderStatus       `json:"status"`
	CreatedAt time.Time         `json:"created_at"`
	Payload   map[string]string `json:"payload,omitempty"`
}

func NewOrderEvent(orderID int, status OrderStatus, opts ...OrderEventOption) (Event, error) {
	orderEv := &OrderEvent{
		OrderID:   orderID,
		Status:    status,
		CreatedAt: time.Now(),
		Payload:   make(map[string]string),
	}
	for _, opt := range opts {
		opt(orderEv)
	}

	value, err := json.Marshal(orderEv)
	if err != nil {
		return Event{}, fmt.Errorf("json.Marshal: %w", err)
	}

	return Event{
		Topic: "loms.order-events",
		Key:   []byte(strconv.Itoa(orderID)),
		Value: value,
	}, nil
}

func OrderEventWithError(err error) OrderEventOption {
	return func(e *OrderEvent) {
		if err != nil {
			e.Payload["error"] = err.Error()
		}
	}
}
