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
	// =========================================================================
	// TODO [Hands-on Step 2.1]: Orkestrasi Domain Constructor & Simpan Produk
	// =========================================================================
	// Sesuai Clean Architecture (Use Case):
	// 1. Panggil constructor entity.NewProduct(req.Name, req.Price, req.Stock):
	// 2. Simpan produk ke repo via interface contracts.ProductRepository:
	// 3. Kembalikan response DTO:

	return nil, fmt.Errorf("TODO [Hands-on Step 2.1]: implementasikan service CreateProduct")
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

func (s *productService) UpdateProduct(ctx context.Context, id string, req dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	// =========================================================================
	// TODO [Hands-on Step 2.2]: Use Case Update Produk
	// =========================================================================
	// 1. Ambil data produk berdasarkan ID:
	// 2. Validasi invariant & update field produk via domain entity method:
	// 3. Simpan perubahan ke repository:
	// 4. Return response:

	return nil, fmt.Errorf("TODO [Hands-on Step 2.2]: implementasikan service UpdateProduct")
}

func (s *productService) DeleteProduct(ctx context.Context, id string) error {
	// =========================================================================
	// TODO [Hands-on Step 2.3]: Use Case Delete Produk
	// =========================================================================
	// 1. Pastikan produk ada di database:
	// 2. Hapus produk via repository:
	// 3. Return nil jika sukses

	return fmt.Errorf("TODO [Hands-on Step 2.3]: implementasikan service DeleteProduct")
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
