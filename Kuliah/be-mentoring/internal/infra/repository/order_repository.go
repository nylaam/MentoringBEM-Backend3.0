package repository

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"github.com/mentoring/toko-backend/domain/contracts"
	"github.com/mentoring/toko-backend/domain/entity"
	"github.com/mentoring/toko-backend/domain/enums"
	"github.com/mentoring/toko-backend/domain/errx"
)

type orderRepository struct {
	db      *sql.DB
	mu      sync.RWMutex
	inMemDB map[string]*entity.Order
}

func NewOrderRepository(db *sql.DB) contracts.OrderRepository {
	return &orderRepository{
		db:      db,
		inMemDB: make(map[string]*entity.Order),
	}
}

func (r *orderRepository) Save(ctx context.Context, o *entity.Order) error {
	if r.db == nil {
		r.mu.Lock()
		defer r.mu.Unlock()
		// Simpan copy ke in-memory map
		cp := *o
		cpItems := make([]entity.OrderItem, len(o.Items))
		copy(cpItems, o.Items)
		cp.Items = cpItems
		r.inMemDB[o.ID] = &cp
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("gagal memulai database transaction: %w", err)
	}
	defer tx.Rollback()

	orderQuery := `INSERT INTO orders (id, customer_id, customer_name, total, status, created_at, updated_at) 
	               VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err = tx.ExecContext(ctx, orderQuery, o.ID, o.CustomerID, o.CustomerName, o.Total, string(o.Status), o.CreatedAt, o.UpdatedAt)
	if err != nil {
		return fmt.Errorf("gagal menyimpan header order: %w", err)
	}

	itemQuery := `INSERT INTO order_items (order_id, product_id, product_name, quantity, unit_price, subtotal) 
	              VALUES ($1, $2, $3, $4, $5, $6)`
	for _, it := range o.Items {
		_, err = tx.ExecContext(ctx, itemQuery, o.ID, it.ProductID, it.ProductName, it.Quantity, it.UnitPrice, it.Subtotal)
		if err != nil {
			return fmt.Errorf("gagal menyimpan order item: %w", err)
		}
	}

	return tx.Commit()
}

func (r *orderRepository) FindByID(ctx context.Context, id string) (*entity.Order, error) {
	if r.db == nil {
		r.mu.RLock()
		defer r.mu.RUnlock()
		o, exists := r.inMemDB[id]
		if !exists {
			return nil, errx.ErrOrderNotFound
		}
		cp := *o
		return &cp, nil
	}

	orderQuery := `SELECT id, customer_id, customer_name, total, status, created_at, updated_at FROM orders WHERE id = $1`
	var o entity.Order
	var statusStr string
	err := r.db.QueryRowContext(ctx, orderQuery, id).Scan(&o.ID, &o.CustomerID, &o.CustomerName, &o.Total, &statusStr, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errx.ErrOrderNotFound
		}
		return nil, fmt.Errorf("gagal query order by id: %w", err)
	}
	o.Status = enums.OrderStatus(statusStr)

	itemsQuery := `SELECT product_id, product_name, quantity, unit_price, subtotal FROM order_items WHERE order_id = $1`
	rows, err := r.db.QueryContext(ctx, itemsQuery, id)
	if err != nil {
		return nil, fmt.Errorf("gagal query items order: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var it entity.OrderItem
		if err := rows.Scan(&it.ProductID, &it.ProductName, &it.Quantity, &it.UnitPrice, &it.Subtotal); err != nil {
			return nil, err
		}
		o.Items = append(o.Items, it)
	}

	return &o, nil
}

func (r *orderRepository) FindAll(ctx context.Context) ([]*entity.Order, error) {
	if r.db == nil {
		r.mu.RLock()
		defer r.mu.RUnlock()
		var list []*entity.Order
		for _, o := range r.inMemDB {
			cp := *o
			list = append(list, &cp)
		}
		return list, nil
	}

	query := `SELECT id, customer_id, customer_name, total, status, created_at, updated_at FROM orders ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("gagal query list orders: %w", err)
	}
	defer rows.Close()

	var list []*entity.Order
	for rows.Next() {
		var o entity.Order
		var statusStr string
		if err := rows.Scan(&o.ID, &o.CustomerID, &o.CustomerName, &o.Total, &statusStr, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		o.Status = enums.OrderStatus(statusStr)
		list = append(list, &o)
	}
	return list, nil
}
