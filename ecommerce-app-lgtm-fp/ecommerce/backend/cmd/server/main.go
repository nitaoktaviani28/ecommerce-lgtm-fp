package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/lgtm-fp/ecommerce-backend/internal/handlers"
	"github.com/lgtm-fp/ecommerce-backend/internal/middleware"
	"github.com/lgtm-fp/ecommerce-backend/internal/repository"
	"github.com/lgtm-fp/ecommerce-backend/observability"
)

func main() {
	// ─── 1. Inisialisasi Observability ───────────────────────────────────────
	// Single entry point — tracing (Tempo), profiling (Pyroscope), metrics (Mimir).
	// Tidak ada perubahan pada logika bisnis untuk mengaktifkan observability.
	observability.Init()

	// ─── 2. Koneksi Database dengan Tracing ──────────────────────────────────
	// OpenDBWithTracing membuka koneksi PostgreSQL yang diinstrumentasi otelsql.
	// Setiap query SQL akan menghasilkan span child di Tempo secara otomatis.
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		observability.GetEnv("DB_HOST", "postgres.ecommerce.svc.cluster.local"),
		observability.GetEnv("DB_PORT", "5432"),
		observability.GetEnv("DB_USER", "ecommerce"),
		observability.GetEnv("DB_PASSWORD", "secret"),
		observability.GetEnv("DB_NAME", "ecommerce"),
		observability.GetEnv("DB_SSLMODE", "disable"),
	)

	db, err := observability.OpenDBWithTracing(dsn)
	if err != nil {
		log.Fatalf("Gagal koneksi ke database: %v", err)
	}
	defer db.Close()

	// Verifikasi koneksi DB
	if err := db.Ping(); err != nil {
		log.Fatalf("Database tidak dapat dijangkau: %v", err)
	}
	log.Println("Database connected with tracing enabled")

	// ─── 3. Inisialisasi Repository ──────────────────────────────────────────
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)

	// ─── 4. Inisialisasi Handler ─────────────────────────────────────────────
	productHandler := handlers.NewProductHandler(productRepo)
	orderHandler := handlers.NewOrderHandler(orderRepo, productRepo)

	// ─── 5. Setup Router ─────────────────────────────────────────────────────
	mux := http.NewServeMux()

	// Health check — digunakan oleh Kubernetes liveness & readiness probe
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Metrics endpoint — di-scrape oleh Alloy untuk dikirim ke Mimir
	mux.Handle("/metrics", promhttp.Handler())

	// Product endpoints
	mux.HandleFunc("/api/products", productHandler.GetProducts)
	mux.HandleFunc("/api/products/", productHandler.GetProduct)

	// Order endpoints
	mux.HandleFunc("/api/orders", orderHandler.CreateOrder)
	mux.HandleFunc("/api/orders/", orderHandler.GetOrder)

	// User orders
	mux.HandleFunc("/api/users/", orderHandler.GetUserOrders)

	// ─── 6. Apply Middleware ──────────────────────────────────────────────────
	// Telemetry middleware: OTel tracing + Prometheus metrics + structured logging
	// CORS middleware: izinkan request dari FE Vue.js
	handler := middleware.Telemetry(middleware.CORS(mux))

	// ─── 7. Start Server ─────────────────────────────────────────────────────
	port := observability.GetEnv("PORT", "8080")
	log.Printf("Server running on :%s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
