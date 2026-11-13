package errx

var (
	ErrProductNotFound     = New(KindNotFound, "PRODUCT_NOT_FOUND", "produk tidak ditemukan")
	ErrInsufficientStock   = New(KindInvalid, "INSUFFICIENT_STOCK", "stok produk tidak mencukupi")
	ErrInvalidProductName  = New(KindInvalid, "INVALID_NAME", "nama produk minimal harus 3 karakter")
	ErrInvalidProductPrice = New(KindInvalid, "INVALID_PRICE", "harga produk harus lebih besar dari 0")
	ErrInvalidProductStock = New(KindInvalid, "INVALID_STOCK", "stok produk tidak boleh negatif")
)
