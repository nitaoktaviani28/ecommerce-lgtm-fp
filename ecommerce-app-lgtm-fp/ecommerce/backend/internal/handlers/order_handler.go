package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/lgtm-fp/ecommerce-backend/internal/domain"
	"github.com/lgtm-fp/ecommerce-backend/internal/repository"
)

type OrderHandler struct {
	orderRepo   *repository.OrderRepository
	productRepo *repository.ProductRepository
}

func NewOrderHandler(orderRepo *repository.OrderRepository, productRepo *repository.ProductRepository) *OrderHandler {
	return &OrderHandler{orderRepo: orderRepo, productRepo: productRepo}
}

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

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := strings.TrimPrefix(r.URL.Path, "/api/orders/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID tidak valid")
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
