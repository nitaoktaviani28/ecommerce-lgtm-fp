package domain

import "time"

// Product merepresentasikan entitas produk dalam sistem ecommerce.
type Product struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	Category    string    `json:"category"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Order merepresentasikan entitas order yang dibuat oleh user.
type Order struct {
	ID         int64       `json:"id"`
	UserID     int64       `json:"user_id"`
	Status     string      `json:"status"` // pending, paid, shipped, completed, cancelled
	TotalPrice float64     `json:"total_price"`
	Items      []OrderItem `json:"items"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

// OrderItem merepresentasikan item dalam sebuah order.
type OrderItem struct {
	ID        int64   `json:"id"`
	OrderID   int64   `json:"order_id"`
	ProductID int64   `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

// User merepresentasikan entitas pengguna sistem.
type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateOrderRequest adalah struktur request untuk membuat order baru.
type CreateOrderRequest struct {
	UserID int64             `json:"user_id"`
	Items  []OrderItemRequest `json:"items"`
}

// OrderItemRequest adalah struktur item dalam request order.
type OrderItemRequest struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}
