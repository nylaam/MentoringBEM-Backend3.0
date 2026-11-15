package dto

import "time"

type CreateProductRequest struct {
	Name  string `json:"name" validate:"required,min=3,max=100"`
	Price int64  `json:"price" validate:"required,gt=0"`
	Stock int    `json:"stock" validate:"required,gte=0"`
}

type UpdateProductRequest struct {
	Name  string `json:"name" validate:"required,min=3,max=100"`
	Price int64  `json:"price" validate:"required,gt=0"`
	Stock int    `json:"stock" validate:"required,gte=0"`
}

// PatchProductRequest mendukung update parsial (hanya field non-nil yang diubah).
// Cocok untuk perubahan cepat di Frontend (misalnya update stok setelah restock, ubah harga, ganti nama produk).
type PatchProductRequest struct {
	Name  *string `json:"name,omitempty" validate:"omitempty,min=3,max=100"`
	Price *int64  `json:"price,omitempty" validate:"omitempty,gt=0"`
	Stock *int    `json:"stock,omitempty" validate:"omitempty,gte=0"`
}

type ProductResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Price     int64     `json:"price"`
	Stock     int       `json:"stock"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
