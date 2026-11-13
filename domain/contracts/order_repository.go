package contracts

import (
	"context"

	"github.com/mentoring/toko-backend/domain/entity"
)

type OrderRepository interface {
	Save(ctx context.Context, order *entity.Order) error
	FindByID(ctx context.Context, id string) (*entity.Order, error)
	FindAll(ctx context.Context) ([]*entity.Order, error)
	Update(ctx context.Context, order *entity.Order) error
	Delete(ctx context.Context, id string) error
}
