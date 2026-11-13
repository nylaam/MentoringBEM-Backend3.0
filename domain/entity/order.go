package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/mentoring/toko-backend/domain/enums"
	"github.com/mentoring/toko-backend/domain/errx"
)

// Invariant batas maksimal belanja (Slide 16-17 di index.md)
const MaxOrderTotal int64 = 50_000_000

type OrderItem struct {
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	Quantity    int    `json:"quantity"`
	UnitPrice   int64  `json:"unit_price"`
	Subtotal    int64  `json:"subtotal"`
}

type Order struct {
	ID           string            `json:"id"`
	CustomerID   string            `json:"customer_id"`
	CustomerName string            `json:"customer_name"`
	Items        []OrderItem       `json:"items"`
	Total        int64             `json:"total"`
	Status       enums.OrderStatus `json:"status"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

// NewOrder adalah Domain Entity Constructor yang menjaga invariant aturan bisnis.
func NewOrder(customerID, customerName string, items []OrderItem) (*Order, error) {
	if len(items) == 0 {
		return nil, errx.ErrEmptyOrder
	}

	var calculatedTotal int64
	validatedItems := make([]OrderItem, len(items))

	for i, it := range items {
		if it.Quantity <= 0 {
			return nil, errx.ErrInvalidQuantity
		}
		subtotal := it.UnitPrice * int64(it.Quantity)
		validatedItems[i] = OrderItem{
			ProductID:   it.ProductID,
			ProductName: it.ProductName,
			Quantity:    it.Quantity,
			UnitPrice:   it.UnitPrice,
			Subtotal:    subtotal,
		}
		calculatedTotal += subtotal
	}

	// =========================================================================
	// TODO [Hands-on Step 1.2]: Validasi Domain Invariant Batas Maksimal Order
	// =========================================================================
	// Skenario Lab (Slide 16 & 17 di index.md):
	// Toko memiliki aturan mutlak: Total nilai pesanan (Order) tidak boleh
	// melebihi Rp 50.000.000 (MaxOrderTotal).
	//
	// Tulis logika validasi di sini:
	// Jika calculatedTotal > MaxOrderTotal,
	// kembalikan: nil, errx.ErrOrderLimitExceeded
	//
	// (Lihat kunci jawaban lengkap di answer.md)
	// =========================================================================

	now := time.Now().UTC()
	return &Order{
		ID:           uuid.NewString(),
		CustomerID:   customerID,
		CustomerName: customerName,
		Items:        validatedItems,
		Total:        calculatedTotal,
		Status:       enums.OrderStatusPending,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}
