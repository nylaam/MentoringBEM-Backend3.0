package unit_test

import (
	"testing"
	"time"

	"github.com/mentoring/toko-backend/domain/entity"
	"github.com/mentoring/toko-backend/domain/errx"
	"github.com/stretchr/testify/assert"
)

func TestNewProduct_Validation(t *testing.T) {
	tests := []struct {
		name      string
		prodName  string
		price     int64
		stock     int
		wantErr   error
		wantStock int
	}{
		{
			name:      "Produk Normal Sukses",
			prodName:  "Mechanical Keyboard",
			price:     1_200_000,
			stock:     15,
			wantErr:   nil,
			wantStock: 15,
		},
		{
			name:     "Gagal Jika Nama Kosong",
			prodName: "",
			price:    1_200_000,
			stock:    15,
			wantErr:  errx.ErrInvalidProductName,
		},
		{
			name:     "Gagal Jika Nama Kurang Dari Tiga Karakter",
			prodName: "AB",
			price:    1_200_000,
			stock:    15,
			wantErr:  errx.ErrInvalidProductName,
		},
		{
			name:     "Gagal Jika Harga Nol Atau Negatif",
			prodName: "Mechanical Keyboard",
			price:    0,
			stock:    15,
			wantErr:  errx.ErrInvalidProductPrice,
		},
		{
			name:     "Gagal Jika Stok Negatif",
			prodName: "Mechanical Keyboard",
			price:    1_200_000,
			stock:    -1,
			wantErr:  errx.ErrInvalidProductStock,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := entity.NewProduct(tt.prodName, tt.price, tt.stock)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, p)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, p)
				assert.Equal(t, tt.prodName, p.Name)
				assert.Equal(t, tt.price, p.Price)
				assert.Equal(t, tt.wantStock, p.Stock)
			}
		})
	}
}

func TestProduct_Update(t *testing.T) {
	p := &entity.Product{
		ID:        "prod-1",
		Name:      "Monitor Lama",
		Price:     1_000_000,
		Stock:     5,
		CreatedAt: time.Now().Add(-1 * time.Hour),
		UpdatedAt: time.Now().Add(-1 * time.Hour),
	}

	oldUpdatedAt := p.UpdatedAt

	// 1. Update sukses
	err := p.Update("Monitor Gaming Baru", 2_500_000, 10)
	assert.NoError(t, err)
	assert.Equal(t, "Monitor Gaming Baru", p.Name)
	assert.Equal(t, int64(2_500_000), p.Price)
	assert.Equal(t, 10, p.Stock)
	assert.True(t, p.UpdatedAt.After(oldUpdatedAt))

	// 2. Gagal jika nama tidak valid
	err = p.Update("X", 2_500_000, 10)
	assert.ErrorIs(t, err, errx.ErrInvalidProductName)

	// 3. Gagal jika harga <= 0
	err = p.Update("Monitor Gaming Baru", 0, 10)
	assert.ErrorIs(t, err, errx.ErrInvalidProductPrice)

	// 4. Gagal jika stok < 0
	err = p.Update("Monitor Gaming Baru", 2_500_000, -5)
	assert.ErrorIs(t, err, errx.ErrInvalidProductStock)
}
