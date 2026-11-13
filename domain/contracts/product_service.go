package contracts

import (
	"context"

	"github.com/mentoring/toko-backend/domain/dto"
)

type ProductService interface {
	CreateProduct(ctx context.Context, req dto.CreateProductRequest) (*dto.ProductResponse, error)
	GetAllProducts(ctx context.Context) ([]*dto.ProductResponse, error)
	GetProductByID(ctx context.Context, id string) (*dto.ProductResponse, error)
	UpdateProduct(ctx context.Context, id string, req dto.UpdateProductRequest) (*dto.ProductResponse, error)
	DeleteProduct(ctx context.Context, id string) error
}
