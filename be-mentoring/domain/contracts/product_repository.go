package contracts

import (
	"context"

	"github.com/mentoring/toko-backend/domain/entity"
)

type ProductRepository interface {
	Save(ctx context.Context, product *entity.Product) error
	FindByID(ctx context.Context, id string) (*entity.Product, error)
	FindAll(ctx context.Context) ([]*entity.Product, error)
	UpdateStock(ctx context.Context, id string, quantityChange int) error
}
