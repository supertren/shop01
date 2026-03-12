package models

import "time"

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type Order struct {
	ID        int
	CreatedAt time.Time
	Status    OrderStatus
	Total     float64
	Items     []OrderItem
}

type OrderItem struct {
	ProductID int
	Name      string
	Quantity  int
	Price     float64
}

type CartItem struct {
	Product  Product
	Quantity int
}

func (c CartItem) Subtotal() float64 {
	return c.Product.Price * float64(c.Quantity)
}
