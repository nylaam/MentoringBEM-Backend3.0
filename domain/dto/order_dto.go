package dto

import (
	"time"

	"github.com/mentoring/toko-backend/domain/enums"
)

type CreateOrderItemRequest struct {
	ProductID string `json:"product_id" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required,gt=0"`
}

type CreateOrderRequest struct {
	CustomerID   string                   `json:"customer_id" validate:"required"`
	CustomerName string                   `json:"customer_name" validate:"required"`
	Items        []CreateOrderItemRequest `json:"items" validate:"required,min=1,dive"`
}

type UpdateOrderRequest struct {
	Status       enums.OrderStatus `json:"status" validate:"required"`
	CustomerName string            `json:"customer_name,omitempty"`
}

type OrderItemResponse struct {
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	Quantity    int    `json:"quantity"`
	UnitPrice   int64  `json:"unit_price"`
	Subtotal    int64  `json:"subtotal"`
}

type OrderResponse struct {
	ID           string              `json:"id"`
	CustomerID   string              `json:"customer_id"`
	CustomerName string              `json:"customer_name"`
	Items        []OrderItemResponse `json:"items"`
	Total        int64               `json:"total"`
	Status       enums.OrderStatus   `json:"status"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
}
