package unit_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/mentoring/toko-backend/domain/dto"
	"github.com/mentoring/toko-backend/domain/errx"
	"github.com/mentoring/toko-backend/internal/app/product"
	"github.com/mentoring/toko-backend/internal/middlewares"
	"github.com/stretchr/testify/assert"
)

type mockProductService struct {
	createProductFunc  func(ctx context.Context, req dto.CreateProductRequest) (*dto.ProductResponse, error)
	getAllProductsFunc func(ctx context.Context) ([]*dto.ProductResponse, error)
	getProductByIDFunc func(ctx context.Context, id string) (*dto.ProductResponse, error)
	updateProductFunc  func(ctx context.Context, id string, req dto.UpdateProductRequest) (*dto.ProductResponse, error)
	deleteProductFunc  func(ctx context.Context, id string) error
}

func (m *mockProductService) CreateProduct(ctx context.Context, req dto.CreateProductRequest) (*dto.ProductResponse, error) {
	if m.createProductFunc != nil {
		return m.createProductFunc(ctx, req)
	}
	return nil, nil
}

func (m *mockProductService) GetAllProducts(ctx context.Context) ([]*dto.ProductResponse, error) {
	if m.getAllProductsFunc != nil {
		return m.getAllProductsFunc(ctx)
	}
	return nil, nil
}

func (m *mockProductService) GetProductByID(ctx context.Context, id string) (*dto.ProductResponse, error) {
	if m.getProductByIDFunc != nil {
		return m.getProductByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockProductService) UpdateProduct(ctx context.Context, id string, req dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	if m.updateProductFunc != nil {
		return m.updateProductFunc(ctx, id, req)
	}
	return nil, nil
}

func (m *mockProductService) DeleteProduct(ctx context.Context, id string) error {
	if m.deleteProductFunc != nil {
		return m.deleteProductFunc(ctx, id)
	}
	return nil
}

func setupProductTestApp(svc *mockProductService) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: middlewares.CustomErrorHandler,
	})
	handler := product.NewProductHandler(svc)
	product.RegisterRoutes(app.Group("/api/v1"), handler)
	return app
}

// -----------------------------------------------------------------------------
// Tests for POST /api/v1/products (Step 3.1)
// -----------------------------------------------------------------------------

func TestProductHandler_Create_Success(t *testing.T) {
	mockSvc := &mockProductService{
		createProductFunc: func(ctx context.Context, req dto.CreateProductRequest) (*dto.ProductResponse, error) {
			return &dto.ProductResponse{
				ID:        "prod-100",
				Name:      req.Name,
				Price:     req.Price,
				Stock:     req.Stock,
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			}, nil
		},
	}

	app := setupProductTestApp(mockSvc)

	payload := dto.CreateProductRequest{
		Name:  "Mechanical Keyboard RGB",
		Price: 1_200_000,
		Stock: 10,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
}

func TestProductHandler_Create_InvalidJSON(t *testing.T) {
	app := setupProductTestApp(&mockProductService{})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewReader([]byte("{invalid-json")))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestProductHandler_Create_ValidationError(t *testing.T) {
	app := setupProductTestApp(&mockProductService{})

	// Payload kosong tidak memenuhi validasi struct
	payload := dto.CreateProductRequest{}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
}

// -----------------------------------------------------------------------------
// Tests for PUT /api/v1/products/:id (Step 3.2)
// -----------------------------------------------------------------------------

func TestProductHandler_Update_Success(t *testing.T) {
	mockSvc := &mockProductService{
		updateProductFunc: func(ctx context.Context, id string, req dto.UpdateProductRequest) (*dto.ProductResponse, error) {
			return &dto.ProductResponse{
				ID:        id,
				Name:      req.Name,
				Price:     req.Price,
				Stock:     req.Stock,
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			}, nil
		},
	}

	app := setupProductTestApp(mockSvc)

	payload := dto.UpdateProductRequest{
		Name:  "Monitor Ultra Sharp 27 Inch",
		Price: 5_000_000,
		Stock: 8,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/products/prod-001", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestProductHandler_Update_InvalidJSON(t *testing.T) {
	app := setupProductTestApp(&mockProductService{})

	req := httptest.NewRequest(http.MethodPut, "/api/v1/products/prod-001", bytes.NewReader([]byte("{broken")))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestProductHandler_Update_ValidationError(t *testing.T) {
	app := setupProductTestApp(&mockProductService{})

	payload := dto.UpdateProductRequest{
		Name: "A", // Kurang dari min=3
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/products/prod-001", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
}

// -----------------------------------------------------------------------------
// Tests for DELETE /api/v1/products/:id (Step 3.3)
// -----------------------------------------------------------------------------

func TestProductHandler_Delete_Success(t *testing.T) {
	mockSvc := &mockProductService{
		deleteProductFunc: func(ctx context.Context, id string) error {
			return nil
		},
	}

	app := setupProductTestApp(mockSvc)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/products/prod-001", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestProductHandler_Delete_NotFound(t *testing.T) {
	mockSvc := &mockProductService{
		deleteProductFunc: func(ctx context.Context, id string) error {
			return errx.ErrProductNotFound
		},
	}

	app := setupProductTestApp(mockSvc)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/products/nonexistent", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}
