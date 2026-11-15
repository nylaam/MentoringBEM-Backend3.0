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
	if len(name) < 3 {
		return nil, errx.ErrInvalidProductName
	}
	if price <= 0 {
		return nil, errx.ErrInvalidProductPrice
	}
	if stock < 0 {
		return nil, errx.ErrInvalidProductStock
	}

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

// Update memvalidasi dan memperbarui seluruh data produk (PUT).
func (p *Product) Update(name string, price int64, stock int) error {
	if len(name) < 3 {
		return errx.ErrInvalidProductName
	}
	if price <= 0 {
		return errx.ErrInvalidProductPrice
	}
	if stock < 0 {
		return errx.ErrInvalidProductStock
	}

	p.Name = name
	p.Price = price
	p.Stock = stock
	p.UpdatedAt = time.Now().UTC()
	return nil
}

// Patch memvalidasi dan memperbarui data produk secara parsial (PATCH).
// Hanya field non-nil yang akan diperbarui (umumnya untuk perubahan kecil seperti update stok atau harga).
func (p *Product) Patch(name *string, price *int64, stock *int) error {
	if name == nil && price == nil && stock == nil {
		return errx.ErrNoFieldsToUpdate
	}

	if name != nil {
		if len(*name) < 3 {
			return errx.ErrInvalidProductName
		}
		p.Name = *name
	}

	if price != nil {
		if *price <= 0 {
			return errx.ErrInvalidProductPrice
		}
		p.Price = *price
	}

	if stock != nil {
		if *stock < 0 {
			return errx.ErrInvalidProductStock
		}
		p.Stock = *stock
	}

	p.UpdatedAt = time.Now().UTC()
	return nil
}
