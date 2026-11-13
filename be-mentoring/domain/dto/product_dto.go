package dto

import "time"

type CreateProductRequest struct {
	Name  string `json:"name" validate:"required,min=3,max=100"`
	Price int64  `json:"price" validate:"required,gt=0"`
	Stock int    `json:"stock" validate:"required,gte=0"`
}

type ProductResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Price     int64     `json:"price"`
	Stock     int       `json:"stock"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
