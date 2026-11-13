package errx

var (
	ErrOrderNotFound   = New(KindNotFound, "ORDER_NOT_FOUND", "order tidak ditemukan")
	ErrEmptyOrder      = New(KindInvalid, "EMPTY_ORDER", "order tidak boleh kosong")
	ErrInvalidQuantity = New(KindInvalid, "INVALID_QUANTITY", "quantity produk harus lebih besar dari 0")

	// =========================================================================
	// TODO [Hands-on Step 1.1]: Definisikan Domain Error Batas Maksimal Order
	// =========================================================================
	// Sesuai materi Slide 16 & 17 di index.md:
	// Buat error domain baru bernama ErrOrderLimitExceeded:
	// - Kind: KindInvalid
	// - Code: "LIMIT_EXCEEDED"
	// - Message: "total order melebihi batas maksimal"
	//
	// (Lihat pembahasan lengkap di answer.md)
	// =========================================================================
	ErrOrderLimitExceeded = New(KindInvalid, "LIMIT_EXCEEDED", "total order melebihi batas maksimal")
)
