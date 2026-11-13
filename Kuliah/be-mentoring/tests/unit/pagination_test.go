package unit_test

import (
	"testing"

	"github.com/mentoring/toko-backend/helpers/pagination"
	"github.com/stretchr/testify/assert"
)

func TestPagination_Helper(t *testing.T) {
	t.Run("Kalkulasi Pagination Normal", func(t *testing.T) {
		p := pagination.New(2, 10, 45)
		assert.Equal(t, 2, p.Page)
		assert.Equal(t, 10, p.Limit)
		assert.Equal(t, int64(45), p.TotalRows)
		assert.Equal(t, 5, p.TotalPages)
		assert.Equal(t, 10, p.Offset())
	})

	t.Run("Default Nilai Jika Input Negatif atau Nol", func(t *testing.T) {
		p := pagination.New(0, 0, 0)
		assert.Equal(t, 1, p.Page)
		assert.Equal(t, 10, p.Limit)
		assert.Equal(t, 1, p.TotalPages)
		assert.Equal(t, 0, p.Offset())
	})

	t.Run("Batas Maksimal Limit (Max 100)", func(t *testing.T) {
		p := pagination.New(1, 200, 50)
		assert.Equal(t, 100, p.Limit)
	})
}
