package errors

import "github.com/mentoring/toko-backend/domain/errx"

var (
	ErrProductNotFound     = errx.ErrProductNotFound
	ErrInsufficientStock   = errx.ErrInsufficientStock
	ErrInvalidProductName  = errx.ErrInvalidProductName
	ErrInvalidProductPrice = errx.ErrInvalidProductPrice
	ErrInvalidProductStock = errx.ErrInvalidProductStock
)
