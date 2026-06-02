package repository

import (
	"context"
	"database/sql"

	"github.com/lgtm-fp/ecommerce-backend/internal/domain"
)

// OrderRepository mengelola operasi database untuk entitas Order.
type OrderRepository struct {
	db *sql.DB
}

// NewOrderRepository membuat instance OrderRepository baru.
func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// Create membuat order baru beserta item-itemnya dalam satu transaksi.
// Transaksi memastikan order dan semua item tersimpan secara atomic.
// Setiap query dalam transaksi ini akan muncul sebagai child span di Tempo.
func (r *OrderRepository) Create(ctx context.Context, req domain.CreateOrderRequest, totalPrice float64) (*domain.Order, error) {
	// Memulai transaksi DB — span "db.begin_tx" akan muncul di Tempo
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Insert order — span "db.exec: INSERT orders" di Tempo
	var order domain.Order
	err = tx.QueryRowContext(ctx, `
		INSERT INTO orders (user_id, status, total_price, created_at, updated_at)
		VALUES ($1, 'pending', $2, NOW(), NOW())
		RETURNING id, user_id, status, total_price, created_at, updated_at
	`, req.UserID, totalPrice).Scan(
		&order.ID, &order.UserID, &order.Status,
		&order.TotalPrice, &order.CreatedAt, &order.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Insert setiap item order — multiple spans di Tempo
	for _, item := range req.Items {
		var price float64
		err = tx.QueryRowContext(ctx, `
			SELECT price FROM products WHERE id = $1
		`, item.ProductID).Scan(&price)
		if err != nil {
			return nil, err
		}

		var orderItem domain.OrderItem
		err = tx.QueryRowContext(ctx, `
			INSERT INTO order_items (order_id, product_id, quantity, price)
			VALUES ($1, $2, $3, $4)
			RETURNING id, order_id, product_id, quantity, price
		`, order.ID, item.ProductID, item.Quantity, price).Scan(
			&orderItem.ID, &orderItem.OrderID,
			&orderItem.ProductID, &orderItem.Quantity, &orderItem.Price,
		)
		if err != nil {
			return nil, err
		}
		order.Items = append(order.Items, orderItem)

		// Kurangi stok produk
		_, err = tx.ExecContext(ctx, `
			UPDATE products
			SET stock = stock - $1, updated_at = NOW()
			WHERE id = $2 AND stock >= $1
		`, item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}
	}

	// Commit transaksi — span "db.commit" di Tempo
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &order, nil
}

// GetByID mengambil order beserta items berdasarkan ID.
func (r *OrderRepository) GetByID(ctx context.Context, id int64) (*domain.Order, error) {
	var order domain.Order
	err := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, status, total_price, created_at, updated_at
		FROM orders WHERE id = $1
	`, id).Scan(
		&order.ID, &order.UserID, &order.Status,
		&order.TotalPrice, &order.CreatedAt, &order.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Ambil items — span terpisah di Tempo
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, order_id, product_id, quantity, price
		FROM order_items WHERE order_id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item domain.OrderItem
		if err := rows.Scan(
			&item.ID, &item.OrderID,
			&item.ProductID, &item.Quantity, &item.Price,
		); err != nil {
			return nil, err
		}
		order.Items = append(order.Items, item)
	}

	return &order, rows.Err()
}

// GetByUserID mengambil semua order milik user tertentu.
func (r *OrderRepository) GetByUserID(ctx context.Context, userID int64) ([]domain.Order, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, status, total_price, created_at, updated_at
		FROM orders WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var o domain.Order
		if err := rows.Scan(
			&o.ID, &o.UserID, &o.Status,
			&o.TotalPrice, &o.CreatedAt, &o.UpdatedAt,
		); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, rows.Err()
}
