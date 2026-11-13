package contracts

import (
	"context"

	"github.com/mentoring/toko-backend/domain/dto"
)

type OrderService interface {
	CreateOrder(ctx context.Context, req dto.CreateOrderRequest) (*dto.OrderResponse, error)
	GetOrderByID(ctx context.Context, id string) (*dto.OrderResponse, error)
	GetAllOrders(ctx context.Context) ([]*dto.OrderResponse, error)
}
