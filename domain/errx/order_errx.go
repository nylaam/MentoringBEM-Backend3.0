package errx

var (
	ErrOrderNotFound      = New(KindNotFound, "ORDER_NOT_FOUND", "order tidak ditemukan")
	ErrEmptyOrder         = New(KindInvalid, "EMPTY_ORDER", "order tidak boleh kosong")
	ErrInvalidQuantity    = New(KindInvalid, "INVALID_QUANTITY", "quantity produk harus lebih besar dari 0")
	ErrOrderLimitExceeded = New(KindInvalid, "LIMIT_EXCEEDED", "total order melebihi batas maksimal")
	ErrInvalidStatus      = New(KindInvalid, "INVALID_STATUS", "status pesanan tidak valid")
)
