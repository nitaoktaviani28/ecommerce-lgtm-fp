package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/lgtm-fp/ecommerce-backend/internal/domain"
	"github.com/lgtm-fp/ecommerce-backend/internal/repository"
)

// OrderHandler mengelola HTTP handler untuk endpoint order.
type OrderHandler struct {
	orderRepo   *repository.OrderRepository
	productRepo *repository.ProductRepository
}

// NewOrderHandler membuat instance OrderHandler baru.
func NewOrderHandler(
	orderRepo *repository.OrderRepository,
	productRepo *repository.ProductRepository,
) *OrderHandler {
	return &OrderHandler{
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

// CreateOrder membuat order baru.
// Flow ini menghasilkan trace lengkap di Tempo:
//   HTTP POST /api/orders
//     └── db.query: SELECT price FROM products (per item)
//     └── db.exec: INSERT INTO orders
//     └── db.exec: INSERT INTO order_items (per item)
//     └── db.exec: UPDATE products SET stock
//     └── db.commit
//
// POST /api/orders
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req domain.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Request tidak valid")
		return
	}

	if len(req.Items) == 0 {
		writeError(w, http.StatusBadRequest, "Order harus memiliki minimal satu item")
		return
	}

	// Hitung total harga dengan mengambil harga tiap produk dari DB.
	// Setiap SELECT ini akan muncul sebagai child span di Tempo.
	var totalPrice float64
	for _, item := range req.Items {
		product, err := h.productRepo.GetByID(ctx, item.ProductID)
		if err != nil || product == nil {
			writeError(w, http.StatusBadRequest, "Produk tidak ditemukan")
			return
		}
		if product.Stock < item.Quantity {
			writeError(w, http.StatusBadRequest, "Stok tidak mencukupi")
			return
		}
		totalPrice += product.Price * float64(item.Quantity)
	}

	order, err := h.orderRepo.Create(ctx, req, totalPrice)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Gagal membuat order")
		return
	}

	writeJSON(w, http.StatusCreated, order)
}

// GetOrder mengambil detail order berdasarkan ID.
//
// GET /api/orders/{id}
func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := strings.TrimPrefix(r.URL.Path, "/api/orders/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID order tidak valid")
		return
	}

	order, err := h.orderRepo.GetByID(ctx, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Gagal mengambil order")
		return
	}
	if order == nil {
		writeError(w, http.StatusNotFound, "Order tidak ditemukan")
		return
	}

	writeJSON(w, http.StatusOK, order)
}

// GetUserOrders mengambil semua order milik user tertentu.
//
// GET /api/users/{user_id}/orders
func (h *OrderHandler) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		writeError(w, http.StatusBadRequest, "User ID tidak valid")
		return
	}

	userID, err := strconv.ParseInt(parts[3], 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "User ID tidak valid")
		return
	}

	orders, err := h.orderRepo.GetByUserID(ctx, userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Gagal mengambil order")
		return
	}

	writeJSON(w, http.StatusOK, orders)
}
