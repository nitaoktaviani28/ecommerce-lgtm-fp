package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/lgtm-fp/ecommerce-backend/internal/repository"
)

// ProductHandler mengelola HTTP handler untuk endpoint produk.
type ProductHandler struct {
	repo *repository.ProductRepository
}

// NewProductHandler membuat instance ProductHandler baru.
func NewProductHandler(repo *repository.ProductRepository) *ProductHandler {
	return &ProductHandler{repo: repo}
}

// GetProducts mengembalikan daftar semua produk.
// Context dari request otomatis membawa trace span,
// sehingga query DB yang dipanggil di bawahnya terhubung ke trace ini.
//
// GET /api/products
// GET /api/products?category=electronics
func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	category := r.URL.Query().Get("category")

	var products interface{}
	var err error

	if category != "" {
		products, err = h.repo.GetByCategory(ctx, category)
	} else {
		products, err = h.repo.GetAll(ctx)
	}

	if err != nil {
		writeError(w, http.StatusInternalServerError, "Gagal mengambil data produk")
		return
	}

	writeJSON(w, http.StatusOK, products)
}

// GetProduct mengembalikan satu produk berdasarkan ID.
//
// GET /api/products/{id}
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID produk tidak valid")
		return
	}

	product, err := h.repo.GetByID(ctx, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Gagal mengambil data produk")
		return
	}
	if product == nil {
		writeError(w, http.StatusNotFound, "Produk tidak ditemukan")
		return
	}

	writeJSON(w, http.StatusOK, product)
}

// writeJSON adalah helper untuk menulis response JSON.
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeError adalah helper untuk menulis response error dalam format JSON.
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
