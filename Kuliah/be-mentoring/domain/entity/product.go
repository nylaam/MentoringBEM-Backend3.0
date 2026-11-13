package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Price     int64     `json:"price"`
	Stock     int       `json:"stock"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewProduct(name string, price int64, stock int) (*Product, error) {
	if name == "" {
		return nil, errors.New("nama produk tidak boleh kosong")
	}
	if price < 0 {
		return nil, errors.New("harga produk tidak boleh negatif")
	}
	if stock < 0 {
		return nil, errors.New("stok produk tidak boleh negatif")
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
