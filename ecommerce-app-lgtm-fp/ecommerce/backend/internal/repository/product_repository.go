package repository

import (
	"context"
	"database/sql"

	"github.com/lgtm-fp/ecommerce-backend/internal/domain"
)

// ProductRepository mengelola operasi database untuk entitas Product.
// Setiap method menerima context yang membawa trace span dari layer handler,
// sehingga query SQL otomatis terhubung ke distributed trace di Tempo.
type ProductRepository struct {
	db *sql.DB
}

// NewProductRepository membuat instance ProductRepository baru.
// db yang diterima harus sudah dibuka menggunakan observability.OpenDBWithTracing()
// agar setiap query otomatis menghasilkan span tracing.
func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// GetAll mengambil seluruh produk dari database.
// Context yang diteruskan memastikan span DB query terhubung
// ke parent span HTTP request di Tempo.
func (r *ProductRepository) GetAll(ctx context.Context) ([]domain.Product, error) {
	// Query ini akan otomatis menghasilkan span "db.query" di Tempo
	// karena menggunakan koneksi yang diinstrumentasi oleh otelsql.
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, description, price, stock, category, created_at, updated_at
		FROM products
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var p domain.Product
		if err := rows.Scan(
			&p.ID, &p.Name, &p.Description,
			&p.Price, &p.Stock, &p.Category,
			&p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, rows.Err()
}

// GetByID mengambil satu produk berdasarkan ID.
func (r *ProductRepository) GetByID(ctx context.Context, id int64) (*domain.Product, error) {
	var p domain.Product
	err := r.db.QueryRowContext(ctx, `
		SELECT id, name, description, price, stock, category, created_at, updated_at
		FROM products
		WHERE id = $1
	`, id).Scan(
		&p.ID, &p.Name, &p.Description,
		&p.Price, &p.Stock, &p.Category,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &p, err
}

// GetByCategory mengambil produk berdasarkan kategori.
func (r *ProductRepository) GetByCategory(ctx context.Context, category string) ([]domain.Product, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, description, price, stock, category, created_at, updated_at
		FROM products
		WHERE category = $1
		ORDER BY created_at DESC
	`, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var p domain.Product
		if err := rows.Scan(
			&p.ID, &p.Name, &p.Description,
			&p.Price, &p.Stock, &p.Category,
			&p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, rows.Err()
}

// DecrementStock mengurangi stok produk secara atomic.
// Menggunakan FOR UPDATE untuk mencegah race condition saat checkout bersamaan.
func (r *ProductRepository) DecrementStock(ctx context.Context, tx *sql.Tx, productID int64, quantity int) error {
	_, err := tx.ExecContext(ctx, `
		UPDATE products
		SET stock = stock - $1, updated_at = NOW()
		WHERE id = $2 AND stock >= $1
	`, quantity, productID)
	return err
}
