package errors

import "github.com/mentoring/toko-backend/domain/errx"

var (
	ErrOrderNotFound      = errx.ErrOrderNotFound
	ErrEmptyOrder         = errx.ErrEmptyOrder
	ErrInvalidQuantity    = errx.ErrInvalidQuantity
	ErrOrderLimitExceeded = errx.ErrOrderLimitExceeded
	ErrInvalidStatus      = errx.ErrInvalidStatus
)
