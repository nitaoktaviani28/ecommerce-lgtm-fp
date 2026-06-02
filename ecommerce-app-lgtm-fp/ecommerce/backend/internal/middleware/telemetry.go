package middleware

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// Prometheus metrics untuk HTTP layer.
// Metrics ini akan di-scrape oleh Alloy dan dikirim ke Mimir.
var (
	httpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total jumlah HTTP request yang diterima.",
	}, []string{"method", "path", "status"})

	httpRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Durasi HTTP request dalam detik.",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "path"})

	httpRequestsInFlight = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "http_requests_in_flight",
		Help: "Jumlah HTTP request yang sedang diproses.",
	})
)

// responseWriter wrapper untuk menangkap status code HTTP response.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Telemetry menggabungkan tiga middleware dalam satu fungsi:
//  1. OpenTelemetry tracing — setiap request menghasilkan span di Tempo
//  2. Prometheus metrics — counter, histogram, gauge
//  3. Structured logging — log setiap request dengan durasi dan status
//
// Middleware ini dipasang di level router sehingga seluruh endpoint
// otomatis ter-observasi tanpa perubahan pada handler.
func Telemetry(next http.Handler) http.Handler {
	// Wrap handler dengan otelhttp untuk menghasilkan span HTTP otomatis.
	// Span ini menjadi root span untuk setiap request,
	// dan semua operasi DB di bawahnya akan menjadi child span.
	otelHandler := otelhttp.NewHandler(next, "http.request",
		otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents),
	)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := newResponseWriter(w)

		httpRequestsInFlight.Inc()
		defer httpRequestsInFlight.Dec()

		// Jalankan handler dengan tracing aktif
		otelHandler.ServeHTTP(rw, r)

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(rw.statusCode)
		path := r.URL.Path

		// Record Prometheus metrics
		httpRequestsTotal.WithLabelValues(r.Method, path, status).Inc()
		httpRequestDuration.WithLabelValues(r.Method, path).Observe(duration)

		// Structured log setiap request
		log.Printf(
			"method=%s path=%s status=%s duration=%.4fs",
			r.Method, path, status, duration,
		)
	})
}

// CORS middleware untuk mengizinkan request dari FE Vue.js.
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
