package repository

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/mentoring/toko-backend/domain/contracts"
	"github.com/mentoring/toko-backend/domain/entity"
	"github.com/mentoring/toko-backend/domain/errx"
)

type productRepository struct {
	db      *sql.DB
	mu      sync.RWMutex
	inMemDB map[string]*entity.Product
}

func NewProductRepository(db *sql.DB) contracts.ProductRepository {
	repo := &productRepository{
		db:      db,
		inMemDB: make(map[string]*entity.Product),
	}

	// Seed default in-memory data untuk hands-on jika database nil
	if db == nil {
		p1 := &entity.Product{ID: "prod-001", Name: "Laptop Gaming ASUS ROG", Price: 25_000_000, Stock: 10, CreatedAt: time.Now(), UpdatedAt: time.Now()}
		p2 := &entity.Product{ID: "prod-002", Name: "Mechanical Keyboard", Price: 1_200_000, Stock: 50, CreatedAt: time.Now(), UpdatedAt: time.Now()}
		p3 := &entity.Product{ID: "prod-003", Name: "Wireless Mouse Ergonomis", Price: 650_000, Stock: 100, CreatedAt: time.Now(), UpdatedAt: time.Now()}
		p4 := &entity.Product{ID: "prod-004", Name: "Super Workstation Server", Price: 60_000_000, Stock: 5, CreatedAt: time.Now(), UpdatedAt: time.Now()}
		repo.inMemDB[p1.ID] = p1
		repo.inMemDB[p2.ID] = p2
		repo.inMemDB[p3.ID] = p3
		repo.inMemDB[p4.ID] = p4
	}

	return repo
}

func (r *productRepository) Save(ctx context.Context, p *entity.Product) error {
	if r.db == nil {
		r.mu.Lock()
		defer r.mu.Unlock()
		r.inMemDB[p.ID] = p
		return nil
	}

	query := `INSERT INTO products (id, name, price, stock, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, query, p.ID, p.Name, p.Price, p.Stock, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		return fmt.Errorf("gagal menyimpan produk ke database: %w", err)
	}
	return nil
}

func (r *productRepository) FindByID(ctx context.Context, id string) (*entity.Product, error) {
	if r.db == nil {
		r.mu.RLock()
		defer r.mu.RUnlock()
		p, exists := r.inMemDB[id]
		if !exists {
			return nil, errx.ErrProductNotFound
		}
		// Return copy
		cp := *p
		return &cp, nil
	}

	query := `SELECT id, name, price, stock, created_at, updated_at FROM products WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var p entity.Product
	err := row.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errx.ErrProductNotFound
		}
		return nil, fmt.Errorf("gagal query produk by id: %w", err)
	}
	return &p, nil
}

func (r *productRepository) FindAll(ctx context.Context) ([]*entity.Product, error) {
	if r.db == nil {
		r.mu.RLock()
		defer r.mu.RUnlock()
		var list []*entity.Product
		for _, p := range r.inMemDB {
			cp := *p
			list = append(list, &cp)
		}
		return list, nil
	}

	query := `SELECT id, name, price, stock, created_at, updated_at FROM products ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("gagal query daftar produk: %w", err)
	}
	defer rows.Close()

	var list []*entity.Product
	for rows.Next() {
		var p entity.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, &p)
	}
	return list, nil
}

func (r *productRepository) UpdateStock(ctx context.Context, id string, quantityChange int) error {
	if r.db == nil {
		r.mu.Lock()
		defer r.mu.Unlock()
		p, exists := r.inMemDB[id]
		if !exists {
			return errx.ErrProductNotFound
		}
		if p.Stock+quantityChange < 0 {
			return errx.ErrInsufficientStock
		}
		p.Stock += quantityChange
		p.UpdatedAt = time.Now().UTC()
		return nil
	}

	query := `UPDATE products SET stock = stock + $1, updated_at = $2 WHERE id = $3 AND stock + $1 >= 0`
	res, err := r.db.ExecContext(ctx, query, quantityChange, time.Now().UTC(), id)
	if err != nil {
		return fmt.Errorf("gagal update stok produk: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errx.ErrInsufficientStock
	}
	return nil
}

func (r *productRepository) Update(ctx context.Context, p *entity.Product) error {
	if r.db == nil {
		r.mu.Lock()
		defer r.mu.Unlock()
		if _, exists := r.inMemDB[p.ID]; !exists {
			return errx.ErrProductNotFound
		}
		cp := *p
		r.inMemDB[p.ID] = &cp
		return nil
	}

	query := `UPDATE products SET name = $1, price = $2, stock = $3, updated_at = $4 WHERE id = $5`
	res, err := r.db.ExecContext(ctx, query, p.Name, p.Price, p.Stock, p.UpdatedAt, p.ID)
	if err != nil {
		return fmt.Errorf("gagal update produk: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errx.ErrProductNotFound
	}
	return nil
}

func (r *productRepository) Delete(ctx context.Context, id string) error {
	if r.db == nil {
		r.mu.Lock()
		defer r.mu.Unlock()
		if _, exists := r.inMemDB[id]; !exists {
			return errx.ErrProductNotFound
		}
		delete(r.inMemDB, id)
		return nil
	}

	query := `DELETE FROM products WHERE id = $1`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("gagal menghapus produk: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errx.ErrProductNotFound
	}
	return nil
}
