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
	"github.com/mentoring/toko-backend/domain/enums"
	"github.com/mentoring/toko-backend/domain/errx"
	"github.com/mentoring/toko-backend/internal/app/order"
	"github.com/mentoring/toko-backend/internal/middlewares"
	"github.com/stretchr/testify/assert"
)

// =============================================================================
// Mock Service Manual untuk Delivery / Controller Test (Slide 12 di index.md)
// =============================================================================

type mockOrderService struct {
	createOrderFunc  func(ctx context.Context, req dto.CreateOrderRequest) (*dto.OrderResponse, error)
	getOrderByIDFunc func(ctx context.Context, id string) (*dto.OrderResponse, error)
	getAllOrdersFunc func(ctx context.Context) ([]*dto.OrderResponse, error)
	updateOrderFunc  func(ctx context.Context, id string, req dto.UpdateOrderRequest) (*dto.OrderResponse, error)
	deleteOrderFunc  func(ctx context.Context, id string) error
}

func (m *mockOrderService) CreateOrder(ctx context.Context, req dto.CreateOrderRequest) (*dto.OrderResponse, error) {
	if m.createOrderFunc != nil {
		return m.createOrderFunc(ctx, req)
	}
	return nil, nil
}

func (m *mockOrderService) GetOrderByID(ctx context.Context, id string) (*dto.OrderResponse, error) {
	if m.getOrderByIDFunc != nil {
		return m.getOrderByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockOrderService) GetAllOrders(ctx context.Context) ([]*dto.OrderResponse, error) {
	if m.getAllOrdersFunc != nil {
		return m.getAllOrdersFunc(ctx)
	}
	return nil, nil
}

func (m *mockOrderService) UpdateOrder(ctx context.Context, id string, req dto.UpdateOrderRequest) (*dto.OrderResponse, error) {
	if m.updateOrderFunc != nil {
		return m.updateOrderFunc(ctx, id, req)
	}
	return nil, nil
}

func (m *mockOrderService) DeleteOrder(ctx context.Context, id string) error {
	if m.deleteOrderFunc != nil {
		return m.deleteOrderFunc(ctx, id)
	}
	return nil
}

func setupTestApp(svc *mockOrderService) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: middlewares.CustomErrorHandler,
	})
	handler := order.NewOrderHandler(svc)
	order.RegisterRoutes(app.Group("/api/v1"), handler)
	return app
}

// -----------------------------------------------------------------------------
// Tests for POST /api/v1/orders (Step 3.1)
// -----------------------------------------------------------------------------

func TestOrderHandler_Create_Success(t *testing.T) {
	mockSvc := &mockOrderService{
		createOrderFunc: func(ctx context.Context, req dto.CreateOrderRequest) (*dto.OrderResponse, error) {
			return &dto.OrderResponse{
				ID:           "order-123",
				CustomerID:   req.CustomerID,
				CustomerName: req.CustomerName,
				Total:        1_000_000,
				Status:       enums.OrderStatusPending,
				CreatedAt:    time.Now().UTC(),
				UpdatedAt:    time.Now().UTC(),
			}, nil
		},
	}

	app := setupTestApp(mockSvc)

	payload := dto.CreateOrderRequest{
		CustomerID:   "cust-001",
		CustomerName: "Budi Santoso",
		Items: []dto.CreateOrderItemRequest{
			{ProductID: "prod-1", Quantity: 2},
		},
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
}

func TestOrderHandler_Create_InvalidJSON(t *testing.T) {
	app := setupTestApp(&mockOrderService{})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewReader([]byte("{invalid-json")))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestOrderHandler_Create_ValidationError(t *testing.T) {
	app := setupTestApp(&mockOrderService{})

	payload := dto.CreateOrderRequest{}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
}

func TestOrderHandler_Create_DomainError(t *testing.T) {
	mockSvc := &mockOrderService{
		createOrderFunc: func(ctx context.Context, req dto.CreateOrderRequest) (*dto.OrderResponse, error) {
			return nil, errx.ErrOrderLimitExceeded
		},
	}

	app := setupTestApp(mockSvc)

	payload := dto.CreateOrderRequest{
		CustomerID:   "cust-001",
		CustomerName: "Budi Santoso",
		Items: []dto.CreateOrderItemRequest{
			{ProductID: "prod-high", Quantity: 1},
		},
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
}

// -----------------------------------------------------------------------------
// Tests for PUT /api/v1/orders/:id (Step 3.2)
// -----------------------------------------------------------------------------

func TestOrderHandler_Update_Success(t *testing.T) {
	mockSvc := &mockOrderService{
		updateOrderFunc: func(ctx context.Context, id string, req dto.UpdateOrderRequest) (*dto.OrderResponse, error) {
			return &dto.OrderResponse{
				ID:           id,
				CustomerID:   "cust-001",
				CustomerName: "Budi Santoso",
				Total:        1_000_000,
				Status:       req.Status,
				CreatedAt:    time.Now().UTC(),
				UpdatedAt:    time.Now().UTC(),
			}, nil
		},
	}

	app := setupTestApp(mockSvc)

	payload := dto.UpdateOrderRequest{
		Status: enums.OrderStatusPaid,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/orders/order-123", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestOrderHandler_Update_InvalidJSON(t *testing.T) {
	app := setupTestApp(&mockOrderService{})

	req := httptest.NewRequest(http.MethodPut, "/api/v1/orders/order-123", bytes.NewReader([]byte("{invalid-json")))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestOrderHandler_Update_ValidationError(t *testing.T) {
	app := setupTestApp(&mockOrderService{})

	// Status kosong melanggar tag validasi required
	payload := dto.UpdateOrderRequest{}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/orders/order-123", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
}

// -----------------------------------------------------------------------------
// Tests for DELETE /api/v1/orders/:id (Step 3.3)
// -----------------------------------------------------------------------------

func TestOrderHandler_Delete_Success(t *testing.T) {
	mockSvc := &mockOrderService{
		deleteOrderFunc: func(ctx context.Context, id string) error {
			return nil
		},
	}

	app := setupTestApp(mockSvc)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/orders/order-123", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestOrderHandler_Delete_NotFound(t *testing.T) {
	mockSvc := &mockOrderService{
		deleteOrderFunc: func(ctx context.Context, id string) error {
			return errx.ErrOrderNotFound
		},
	}

	app := setupTestApp(mockSvc)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/orders/nonexistent-id", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}
