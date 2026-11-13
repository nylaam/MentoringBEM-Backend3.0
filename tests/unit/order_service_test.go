package unit_test

import (
	"context"
	"errors"
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

type mockProductRepo struct {
	findByIDFunc    func(ctx context.Context, id string) (*entity.Product, error)
	findAllFunc     func(ctx context.Context) ([]*entity.Product, error)
	saveFunc        func(ctx context.Context, p *entity.Product) error
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
			return errors.New("database connection timeout")
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
