package unit_test

import (
	"testing"

	"github.com/mentoring/toko-backend/domain/entity"
	"github.com/mentoring/toko-backend/domain/errx"
	"github.com/stretchr/testify/assert"
)

// TestNewOrder_Validation menguji seluruh domain invariant pada entity Order (Slide 18 di index.md)
func TestNewOrder_Validation(t *testing.T) {
	tests := []struct {
		name          string
		customerID    string
		customerName  string
		items         []entity.OrderItem
		wantErr       error
		expectedTotal int64
	}{
		{
			name:         "Order Normal Sukses",
			customerID:   "cust-001",
			customerName: "Budi Santoso",
			items: []entity.OrderItem{
				{ProductID: "prod-1", ProductName: "Keyboard Mechanical", Quantity: 2, UnitPrice: 500_000},
				{ProductID: "prod-2", ProductName: "Mouse Wireless", Quantity: 1, UnitPrice: 250_000},
			},
			wantErr:       nil,
			expectedTotal: 1_250_000,
		},
		{
			name:         "Gagal Jika Quantity Kurang Dari Satu",
			customerID:   "cust-001",
			customerName: "Budi Santoso",
			items: []entity.OrderItem{
				{ProductID: "prod-1", ProductName: "Keyboard Mechanical", Quantity: 0, UnitPrice: 500_000},
			},
			wantErr: errx.ErrInvalidQuantity,
		},
		{
			name:         "Gagal Jika Melebihi Batas Maksimal (50 Juta)",
			customerID:   "cust-001",
			customerName: "Budi Santoso",
			items: []entity.OrderItem{
				{ProductID: "prod-high", ProductName: "Server Rack", Quantity: 1, UnitPrice: 60_000_000},
			},
			wantErr: errx.ErrOrderLimitExceeded,
		},
		{
			name:         "Gagal Jika Item Order Kosong",
			customerID:   "cust-001",
			customerName: "Budi Santoso",
			items:        []entity.OrderItem{},
			wantErr:      errx.ErrEmptyOrder,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order, err := entity.NewOrder(tt.customerID, tt.customerName, tt.items)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, order)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, order)
				assert.Equal(t, tt.expectedTotal, order.Total)
				assert.NotEmpty(t, order.ID)
			}
		})
	}
}
