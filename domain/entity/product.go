package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/mentoring/toko-backend/domain/errx"
)

type Product struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Price     int64     `json:"price"`
	Stock     int       `json:"stock"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewProduct adalah Domain Entity Constructor untuk produk baru.
func NewProduct(name string, price int64, stock int) (*Product, error) {
	// =========================================================================
	// TODO [Hands-on Step 1.1]: Validasi Domain Invariant Product Baru
	// =========================================================================
	// Aturan Bisnis Produk:
	// 1. Nama produk wajib diisi dan minimal 3 karakter:
	// 2. Harga produk harus lebih besar dari 0:
	// 3. Stok produk tidak boleh negatif (>= 0):

	now := time.Now().UTC()
	return &Product{
		ID:        uuid.NewString(),
		Name:      name,
		Price:     price,
		Stock:     stock,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// Update memvalidasi dan memperbarui data produk.
func (p *Product) Update(name string, price int64, stock int) error {
	// =========================================================================
	// TODO [Hands-on Step 1.2]: Validasi Invariant & Update Field Produk
	// =========================================================================
	// 1. Validasi nama produk minimal 3 karakter:
	// 2. Validasi harga produk > 0:
	// 3. Stok produk tidak boleh negatif (>= 0):
	// 4. Update data produk dan timestamp update:

	return errx.ErrInvalidProductName
}
