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
	patchProductFunc   func(ctx context.Context, id string, req dto.PatchProductRequest) (*dto.ProductResponse, error)
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

func (m *mockProductService) PatchProduct(ctx context.Context, id string, req dto.PatchProductRequest) (*dto.ProductResponse, error) {
	if m.patchProductFunc != nil {
		return m.patchProductFunc(ctx, id, req)
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
// Tests for POST /api/v1/products
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

	payload := dto.CreateProductRequest{}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
}

// -----------------------------------------------------------------------------
// Tests for GET /api/v1/products & GET /api/v1/products/:id
// -----------------------------------------------------------------------------

func TestProductHandler_GetAll_Success(t *testing.T) {
	mockSvc := &mockProductService{
		getAllProductsFunc: func(ctx context.Context) ([]*dto.ProductResponse, error) {
			return []*dto.ProductResponse{
				{ID: "p1", Name: "Mouse", Price: 100000, Stock: 5},
			}, nil
		},
	}

	app := setupProductTestApp(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/products", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestProductHandler_GetByID_Success(t *testing.T) {
	mockSvc := &mockProductService{
		getProductByIDFunc: func(ctx context.Context, id string) (*dto.ProductResponse, error) {
			return &dto.ProductResponse{ID: id, Name: "Mouse", Price: 100000, Stock: 5}, nil
		},
	}

	app := setupProductTestApp(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/products/p1", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

// -----------------------------------------------------------------------------
// Tests for PUT /api/v1/products/:id
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
// Tests for PATCH /api/v1/products/:id
// -----------------------------------------------------------------------------

func TestProductHandler_Patch_Success(t *testing.T) {
	mockSvc := &mockProductService{
		patchProductFunc: func(ctx context.Context, id string, req dto.PatchProductRequest) (*dto.ProductResponse, error) {
			stock := 100
			if req.Stock != nil {
				stock = *req.Stock
			}
			return &dto.ProductResponse{
				ID:        id,
				Name:      "Wireless Keyboard",
				Price:     500_000,
				Stock:     stock,
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			}, nil
		},
	}

	app := setupProductTestApp(mockSvc)

	newStock := 50
	payload := dto.PatchProductRequest{
		Stock: &newStock,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/products/prod-001", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestProductHandler_Patch_InvalidJSON(t *testing.T) {
	app := setupProductTestApp(&mockProductService{})

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/products/prod-001", bytes.NewReader([]byte("{bad-json")))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestProductHandler_Patch_NoFieldsError(t *testing.T) {
	app := setupProductTestApp(&mockProductService{})

	// Payload kosong `{}` tanpa ada field
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/products/prod-001", bytes.NewReader([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestProductHandler_Patch_ValidationError(t *testing.T) {
	app := setupProductTestApp(&mockProductService{})

	shortName := "X" // min=3
	payload := dto.PatchProductRequest{
		Name: &shortName,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/products/prod-001", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
}

func TestProductHandler_Patch_NotFound(t *testing.T) {
	mockSvc := &mockProductService{
		patchProductFunc: func(ctx context.Context, id string, req dto.PatchProductRequest) (*dto.ProductResponse, error) {
			return nil, errx.ErrProductNotFound
		},
	}

	app := setupProductTestApp(mockSvc)

	newPrice := int64(750_000)
	payload := dto.PatchProductRequest{
		Price: &newPrice,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/products/nonexistent", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

// -----------------------------------------------------------------------------
// Tests for DELETE /api/v1/products/:id
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
