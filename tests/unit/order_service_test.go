package unit_test

import (
	"context"
	stderrs "errors"
	"testing"

	"github.com/mentoring/toko-backend/domain/dto"
	"github.com/mentoring/toko-backend/domain/entity"
	"github.com/mentoring/toko-backend/domain/errx"
	"github.com/mentoring/toko-backend/internal/app/order"
	"github.com/stretchr/testify/assert"
)

// =============================================================================
// Mock Repository Manual (Slide 19 di index.md)
// =============================================================================

type mockOrderRepo struct {
	saveFunc     func(ctx context.Context, o *entity.Order) error
	findByIDFunc func(ctx context.Context, id string) (*entity.Order, error)
	findAllFunc  func(ctx context.Context) ([]*entity.Order, error)
	updateFunc   func(ctx context.Context, o *entity.Order) error
	deleteFunc   func(ctx context.Context, id string) error
}

func (m *mockOrderRepo) Save(ctx context.Context, o *entity.Order) error {
	if m.saveFunc != nil {
		return m.saveFunc(ctx, o)
	}
	return nil
}

func (m *mockOrderRepo) FindByID(ctx context.Context, id string) (*entity.Order, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockOrderRepo) FindAll(ctx context.Context) ([]*entity.Order, error) {
	if m.findAllFunc != nil {
		return m.findAllFunc(ctx)
	}
	return nil, nil
}

func (m *mockOrderRepo) Update(ctx context.Context, o *entity.Order) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, o)
	}
	return nil
}

func (m *mockOrderRepo) Delete(ctx context.Context, id string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return nil
}

type mockProductRepo struct {
	findByIDFunc    func(ctx context.Context, id string) (*entity.Product, error)
	findAllFunc     func(ctx context.Context) ([]*entity.Product, error)
	saveFunc        func(ctx context.Context, p *entity.Product) error
	updateFunc      func(ctx context.Context, p *entity.Product) error
	deleteFunc      func(ctx context.Context, id string) error
	updateStockFunc func(ctx context.Context, id string, change int) error
}

func (m *mockProductRepo) FindByID(ctx context.Context, id string) (*entity.Product, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockProductRepo) FindAll(ctx context.Context) ([]*entity.Product, error) {
	if m.findAllFunc != nil {
		return m.findAllFunc(ctx)
	}
	return nil, nil
}

func (m *mockProductRepo) Save(ctx context.Context, p *entity.Product) error {
	if m.saveFunc != nil {
		return m.saveFunc(ctx, p)
	}
	return nil
}

func (m *mockProductRepo) Update(ctx context.Context, p *entity.Product) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, p)
	}
	return nil
}

func (m *mockProductRepo) Delete(ctx context.Context, id string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return nil
}

func (m *mockProductRepo) UpdateStock(ctx context.Context, id string, change int) error {
	if m.updateStockFunc != nil {
		return m.updateStockFunc(ctx, id, change)
	}
	return nil
}

// -----------------------------------------------------------------------------
// Test 1: Skenario Produk Tidak Ditemukan
// -----------------------------------------------------------------------------
func TestOrderService_CreateOrder_ProductNotFound(t *testing.T) {
	mockProdRepo := &mockProductRepo{
		findByIDFunc: func(ctx context.Context, id string) (*entity.Product, error) {
			return nil, errx.ErrProductNotFound
		},
	}
	mockOrderRepo := &mockOrderRepo{}

	svc := order.NewOrderService(mockOrderRepo, mockProdRepo)

	req := dto.CreateOrderRequest{
		CustomerID:   "cust-123",
		CustomerName: "Budi",
		Items: []dto.CreateOrderItemRequest{
			{ProductID: "non-existent-id", Quantity: 1},
		},
	}

	res, err := svc.CreateOrder(context.Background(), req)

	assert.Error(t, err)
	assert.ErrorIs(t, err, errx.ErrProductNotFound)
	assert.Nil(t, res)
}

// -----------------------------------------------------------------------------
// Test 2: Skenario Berhasil Membuat Order (Menguji Orkestrasi Service)
// -----------------------------------------------------------------------------
func TestOrderService_CreateOrder_Success(t *testing.T) {
	mockProdRepo := &mockProductRepo{
		findByIDFunc: func(ctx context.Context, id string) (*entity.Product, error) {
			return &entity.Product{
				ID:    "prod-001",
				Name:  "Mechanical Keyboard",
				Price: 1_200_000,
				Stock: 10,
			}, nil
		},
	}

	mockOrderRepo := &mockOrderRepo{
		saveFunc: func(ctx context.Context, o *entity.Order) error {
			assert.Equal(t, int64(2_400_000), o.Total)
			return nil
		},
	}

	svc := order.NewOrderService(mockOrderRepo, mockProdRepo)

	req := dto.CreateOrderRequest{
		CustomerID:   "cust-001",
		CustomerName: "Budi Santoso",
		Items: []dto.CreateOrderItemRequest{
			{ProductID: "prod-001", Quantity: 2},
		},
	}

	res, err := svc.CreateOrder(context.Background(), req)

	assert.NoError(t, err)
	if assert.NotNil(t, res) {
		assert.Equal(t, int64(2_400_000), res.Total)
		assert.Equal(t, "Budi Santoso", res.CustomerName)
		assert.Len(t, res.Items, 1)
	}
}

// -----------------------------------------------------------------------------
// Test 3: Skenario Database Timeout Saat Menyimpan Order
// -----------------------------------------------------------------------------
func TestOrderService_CreateOrder_RepoError(t *testing.T) {
	mockProdRepo := &mockProductRepo{
		findByIDFunc: func(ctx context.Context, id string) (*entity.Product, error) {
			return &entity.Product{
				ID:    "prod-001",
				Name:  "Mouse Wireless",
				Price: 500_000,
				Stock: 5,
			}, nil
		},
	}

	mockOrderRepo := &mockOrderRepo{
		saveFunc: func(ctx context.Context, o *entity.Order) error {
			return stderrs.New("database connection timeout")
		},
	}

	svc := order.NewOrderService(mockOrderRepo, mockProdRepo)

	req := dto.CreateOrderRequest{
		CustomerID:   "cust-001",
		CustomerName: "Budi Santoso",
		Items: []dto.CreateOrderItemRequest{
			{ProductID: "prod-001", Quantity: 1},
		},
	}

	res, err := svc.CreateOrder(context.Background(), req)

	assert.Error(t, err)
	if err != nil {
		assert.ErrorContains(t, err, "database connection timeout")
	}
	assert.Nil(t, res)
}

// -----------------------------------------------------------------------------
// Test 4: Skenario Melebihi Batas Maksimal Belanja (> Rp 50.000.000)
// -----------------------------------------------------------------------------
func TestOrderService_CreateOrder_ExceedLimit(t *testing.T) {
	mockProdRepo := &mockProductRepo{
		findByIDFunc: func(ctx context.Context, id string) (*entity.Product, error) {
			return &entity.Product{
				ID:    "prod-server",
				Name:  "Enterprise Server Rack",
				Price: 60_000_000,
				Stock: 2,
			}, nil
		},
	}

	mockOrderRepo := &mockOrderRepo{
		saveFunc: func(ctx context.Context, o *entity.Order) error {
			t.Fatal("saveFunc tidak boleh dipanggil jika batas limit dilanggar")
			return nil
		},
	}

	svc := order.NewOrderService(mockOrderRepo, mockProdRepo)

	req := dto.CreateOrderRequest{
		CustomerID:   "cust-001",
		CustomerName: "Budi Santoso",
		Items: []dto.CreateOrderItemRequest{
			{ProductID: "prod-server", Quantity: 1},
		},
	}

	res, err := svc.CreateOrder(context.Background(), req)

	assert.Error(t, err)
	assert.ErrorIs(t, err, errx.ErrOrderLimitExceeded)
	assert.Nil(t, res)
}

// -----------------------------------------------------------------------------
// Test 5: Skenario Update Order
// -----------------------------------------------------------------------------
func TestOrderService_UpdateOrder(t *testing.T) {
	existingOrder, err := entity.NewOrder("cust-001", "Budi", []entity.OrderItem{
		{ProductID: "prod-1", ProductName: "Keyboard", Quantity: 1, UnitPrice: 100_000},
	})
	assert.NoError(t, err)

	mockOrderRepo := &mockOrderRepo{
		findByIDFunc: func(ctx context.Context, id string) (*entity.Order, error) {
			if id == "not-found" {
				return nil, errx.ErrOrderNotFound
			}
			return existingOrder, nil
		},
		updateFunc: func(ctx context.Context, o *entity.Order) error {
			return nil
		},
	}

	svc := order.NewOrderService(mockOrderRepo, &mockProductRepo{})

	// 1. Success update
	req := dto.UpdateOrderRequest{
		Status:       "PAID",
		CustomerName: "Budi Santoso Updated",
	}
	res, err := svc.UpdateOrder(context.Background(), "order-1", req)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "PAID", string(res.Status))
	assert.Equal(t, "Budi Santoso Updated", res.CustomerName)

	// 2. Order not found
	resNotFound, errNotFound := svc.UpdateOrder(context.Background(), "not-found", req)
	assert.ErrorIs(t, errNotFound, errx.ErrOrderNotFound)
	assert.Nil(t, resNotFound)

	// 3. Invalid status
	reqInvalidStatus := dto.UpdateOrderRequest{
		Status: "INVALID_STATUS",
	}
	resInvalid, errInvalid := svc.UpdateOrder(context.Background(), "order-1", reqInvalidStatus)
	assert.ErrorIs(t, errInvalid, errx.ErrInvalidStatus)
	assert.Nil(t, resInvalid)
}

// -----------------------------------------------------------------------------
// Test 6: Skenario Delete Order
// -----------------------------------------------------------------------------
func TestOrderService_DeleteOrder(t *testing.T) {
	mockOrderRepo := &mockOrderRepo{
		findByIDFunc: func(ctx context.Context, id string) (*entity.Order, error) {
			if id == "not-found" {
				return nil, errx.ErrOrderNotFound
			}
			return &entity.Order{ID: id}, nil
		},
		deleteFunc: func(ctx context.Context, id string) error {
			return nil
		},
	}

	svc := order.NewOrderService(mockOrderRepo, &mockProductRepo{})

	// 1. Success delete
	err := svc.DeleteOrder(context.Background(), "order-1")
	assert.NoError(t, err)

	// 2. Delete not found
	err = svc.DeleteOrder(context.Background(), "not-found")
	assert.ErrorIs(t, err, errx.ErrOrderNotFound)
}
