package unit_test

import (
	"context"
	stderrs "errors"
	"testing"
	"time"

	"github.com/mentoring/toko-backend/domain/dto"
	"github.com/mentoring/toko-backend/domain/entity"
	"github.com/mentoring/toko-backend/domain/errx"
	"github.com/mentoring/toko-backend/internal/app/product"
	"github.com/stretchr/testify/assert"
)

// mockProductRepoService untuk pengujian unit ProductService
type mockProductRepoService struct {
	saveFunc        func(ctx context.Context, p *entity.Product) error
	findByIDFunc    func(ctx context.Context, id string) (*entity.Product, error)
	findAllFunc     func(ctx context.Context) ([]*entity.Product, error)
	updateFunc      func(ctx context.Context, p *entity.Product) error
	deleteFunc      func(ctx context.Context, id string) error
	updateStockFunc func(ctx context.Context, id string, quantityChange int) error
}

func (m *mockProductRepoService) Save(ctx context.Context, p *entity.Product) error {
	if m.saveFunc != nil {
		return m.saveFunc(ctx, p)
	}
	return nil
}

func (m *mockProductRepoService) FindByID(ctx context.Context, id string) (*entity.Product, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockProductRepoService) FindAll(ctx context.Context) ([]*entity.Product, error) {
	if m.findAllFunc != nil {
		return m.findAllFunc(ctx)
	}
	return nil, nil
}

func (m *mockProductRepoService) Update(ctx context.Context, p *entity.Product) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, p)
	}
	return nil
}

func (m *mockProductRepoService) Delete(ctx context.Context, id string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return nil
}

func (m *mockProductRepoService) UpdateStock(ctx context.Context, id string, quantityChange int) error {
	if m.updateStockFunc != nil {
		return m.updateStockFunc(ctx, id, quantityChange)
	}
	return nil
}

// -----------------------------------------------------------------------------
// Test CreateProduct
// -----------------------------------------------------------------------------

func TestProductService_CreateProduct_Success(t *testing.T) {
	var savedProduct *entity.Product
	mockRepo := &mockProductRepoService{
		saveFunc: func(ctx context.Context, p *entity.Product) error {
			savedProduct = p
			return nil
		},
	}

	svc := product.NewProductService(mockRepo)

	req := dto.CreateProductRequest{
		Name:  "Mechanical Keyboard RGB",
		Price: 1_250_000,
		Stock: 20,
	}

	res, err := svc.CreateProduct(context.Background(), req)
	assert.NoError(t, err)
	if assert.NotNil(t, res) {
		assert.Equal(t, "Mechanical Keyboard RGB", res.Name)
		assert.Equal(t, int64(1_250_000), res.Price)
		assert.Equal(t, 20, res.Stock)
		assert.NotNil(t, savedProduct)
	}
}

func TestProductService_CreateProduct_ValidationError(t *testing.T) {
	mockRepo := &mockProductRepoService{}
	svc := product.NewProductService(mockRepo)

	req := dto.CreateProductRequest{
		Name:  "A", // Invalid length < 3
		Price: 1_000,
		Stock: 5,
	}

	res, err := svc.CreateProduct(context.Background(), req)
	assert.ErrorIs(t, err, errx.ErrInvalidProductName)
	assert.Nil(t, res)
}

func TestProductService_CreateProduct_RepoError(t *testing.T) {
	mockRepo := &mockProductRepoService{
		saveFunc: func(ctx context.Context, p *entity.Product) error {
			return stderrs.New("database connection failed")
		},
	}

	svc := product.NewProductService(mockRepo)

	req := dto.CreateProductRequest{
		Name:  "Mechanical Keyboard RGB",
		Price: 1_250_000,
		Stock: 20,
	}

	res, err := svc.CreateProduct(context.Background(), req)
	assert.Error(t, err)
	assert.Nil(t, res)
}

// -----------------------------------------------------------------------------
// Test GetAllProducts & GetProductByID
// -----------------------------------------------------------------------------

func TestProductService_GetAllProducts(t *testing.T) {
	mockRepo := &mockProductRepoService{
		findAllFunc: func(ctx context.Context) ([]*entity.Product, error) {
			return []*entity.Product{
				{ID: "p1", Name: "Product 1", Price: 1000, Stock: 10},
				{ID: "p2", Name: "Product 2", Price: 2000, Stock: 20},
			}, nil
		},
	}

	svc := product.NewProductService(mockRepo)

	res, err := svc.GetAllProducts(context.Background())
	assert.NoError(t, err)
	assert.Len(t, res, 2)
}

func TestProductService_GetProductByID_Success(t *testing.T) {
	mockRepo := &mockProductRepoService{
		findByIDFunc: func(ctx context.Context, id string) (*entity.Product, error) {
			return &entity.Product{ID: id, Name: "Found Product", Price: 1000, Stock: 10}, nil
		},
	}

	svc := product.NewProductService(mockRepo)

	res, err := svc.GetProductByID(context.Background(), "p1")
	assert.NoError(t, err)
	if assert.NotNil(t, res) {
		assert.Equal(t, "Found Product", res.Name)
	}
}

// -----------------------------------------------------------------------------
// Test UpdateProduct (PUT)
// -----------------------------------------------------------------------------

func TestProductService_UpdateProduct_Success(t *testing.T) {
	existing := &entity.Product{
		ID:        "prod-001",
		Name:      "Monitor Jadul",
		Price:     2_000_000,
		Stock:     10,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo := &mockProductRepoService{
		findByIDFunc: func(ctx context.Context, id string) (*entity.Product, error) {
			return existing, nil
		},
		updateFunc: func(ctx context.Context, p *entity.Product) error {
			return nil
		},
	}

	svc := product.NewProductService(mockRepo)

	req := dto.UpdateProductRequest{
		Name:  "Monitor Ultrawide",
		Price: 4_500_000,
		Stock: 15,
	}

	res, err := svc.UpdateProduct(context.Background(), "prod-001", req)
	assert.NoError(t, err)
	if assert.NotNil(t, res) {
		assert.Equal(t, "Monitor Ultrawide", res.Name)
		assert.Equal(t, int64(4_500_000), res.Price)
	}
}

func TestProductService_UpdateProduct_NotFound(t *testing.T) {
	mockRepo := &mockProductRepoService{
		findByIDFunc: func(ctx context.Context, id string) (*entity.Product, error) {
			return nil, errx.ErrProductNotFound
		},
	}

	svc := product.NewProductService(mockRepo)

	req := dto.UpdateProductRequest{
		Name:  "Monitor Ultrawide",
		Price: 4_500_000,
		Stock: 15,
	}

	res, err := svc.UpdateProduct(context.Background(), "non-existent", req)
	assert.ErrorIs(t, err, errx.ErrProductNotFound)
	assert.Nil(t, res)
}

// -----------------------------------------------------------------------------
// Test PatchProduct (PATCH)
// -----------------------------------------------------------------------------

func TestProductService_PatchProduct_Success(t *testing.T) {
	existing := &entity.Product{
		ID:        "prod-001",
		Name:      "Wireless Keyboard",
		Price:     500_000,
		Stock:     20,
		CreatedAt: time.Now().Add(-1 * time.Hour),
		UpdatedAt: time.Now().Add(-1 * time.Hour),
	}

	var updatedProduct *entity.Product
	mockRepo := &mockProductRepoService{
		findByIDFunc: func(ctx context.Context, id string) (*entity.Product, error) {
			return existing, nil
		},
		updateFunc: func(ctx context.Context, p *entity.Product) error {
			updatedProduct = p
			return nil
		},
	}

	svc := product.NewProductService(mockRepo)

	newStock := 35
	req := dto.PatchProductRequest{
		Stock: &newStock,
	}

	res, err := svc.PatchProduct(context.Background(), "prod-001", req)
	assert.NoError(t, err)
	if assert.NotNil(t, res) {
		assert.Equal(t, 35, res.Stock)
		assert.Equal(t, "Wireless Keyboard", res.Name)
		assert.Equal(t, int64(500_000), res.Price)
		assert.NotNil(t, updatedProduct)
		assert.Equal(t, 35, updatedProduct.Stock)
	}
}

func TestProductService_PatchProduct_NoFieldsError(t *testing.T) {
	mockRepo := &mockProductRepoService{}
	svc := product.NewProductService(mockRepo)

	req := dto.PatchProductRequest{}

	res, err := svc.PatchProduct(context.Background(), "prod-001", req)
	assert.ErrorIs(t, err, errx.ErrNoFieldsToUpdate)
	assert.Nil(t, res)
}

func TestProductService_PatchProduct_NotFound(t *testing.T) {
	mockRepo := &mockProductRepoService{
		findByIDFunc: func(ctx context.Context, id string) (*entity.Product, error) {
			return nil, errx.ErrProductNotFound
		},
	}

	svc := product.NewProductService(mockRepo)

	newPrice := int64(600_000)
	req := dto.PatchProductRequest{
		Price: &newPrice,
	}

	res, err := svc.PatchProduct(context.Background(), "non-existent", req)
	assert.ErrorIs(t, err, errx.ErrProductNotFound)
	assert.Nil(t, res)
}

// -----------------------------------------------------------------------------
// Test DeleteProduct
// -----------------------------------------------------------------------------

func TestProductService_DeleteProduct_Success(t *testing.T) {
	deleted := false
	mockRepo := &mockProductRepoService{
		findByIDFunc: func(ctx context.Context, id string) (*entity.Product, error) {
			return &entity.Product{ID: id, Name: "Keyboard"}, nil
		},
		deleteFunc: func(ctx context.Context, id string) error {
			deleted = true
			return nil
		},
	}

	svc := product.NewProductService(mockRepo)

	err := svc.DeleteProduct(context.Background(), "prod-001")
	assert.NoError(t, err)
	assert.True(t, deleted)
}

func TestProductService_DeleteProduct_NotFound(t *testing.T) {
	mockRepo := &mockProductRepoService{
		findByIDFunc: func(ctx context.Context, id string) (*entity.Product, error) {
			return nil, errx.ErrProductNotFound
		},
	}

	svc := product.NewProductService(mockRepo)

	err := svc.DeleteProduct(context.Background(), "non-existent")
	assert.ErrorIs(t, err, errx.ErrProductNotFound)
}
