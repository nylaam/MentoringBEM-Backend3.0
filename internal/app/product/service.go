package product

import (
	"context"
	"fmt"

	"github.com/mentoring/toko-backend/domain/contracts"
	"github.com/mentoring/toko-backend/domain/dto"
	"github.com/mentoring/toko-backend/domain/entity"
)

type productService struct {
	repo contracts.ProductRepository
}

func NewProductService(repo contracts.ProductRepository) contracts.ProductService {
	return &productService{repo: repo}
}

func (s *productService) CreateProduct(ctx context.Context, req dto.CreateProductRequest) (*dto.ProductResponse, error) {
	prod, err := entity.NewProduct(req.Name, req.Price, req.Stock)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Save(ctx, prod); err != nil {
		return nil, fmt.Errorf("gagal menyimpan produk: %w", err)
	}

	return toProductResponse(prod), nil
}

func (s *productService) GetAllProducts(ctx context.Context) ([]*dto.ProductResponse, error) {
	products, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var res []*dto.ProductResponse
	for _, p := range products {
		res = append(res, toProductResponse(p))
	}
	return res, nil
}

func (s *productService) GetProductByID(ctx context.Context, id string) (*dto.ProductResponse, error) {
	p, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toProductResponse(p), nil
}

func toProductResponse(p *entity.Product) *dto.ProductResponse {
	return &dto.ProductResponse{
		ID:        p.ID,
		Name:      p.Name,
		Price:     p.Price,
		Stock:     p.Stock,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}
