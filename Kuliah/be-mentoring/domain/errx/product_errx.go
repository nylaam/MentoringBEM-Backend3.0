package errx

var (
	ErrProductNotFound   = New(KindNotFound, "PRODUCT_NOT_FOUND", "produk tidak ditemukan")
	ErrInsufficientStock = New(KindInvalid, "INSUFFICIENT_STOCK", "stok produk tidak mencukupi")
)
