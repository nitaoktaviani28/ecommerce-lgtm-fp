package observability

import (
	"log"
)

// Init berfungsi untuk menginisialisasi seluruh komponen observability.
// Fungsi ini menjadi satu-satunya titik masuk (single entry point)
// yang dipanggil oleh aplikasi utama untuk mengaktifkan tracing,
// profiling, dan metrics.
//
// Pendekatan ini memastikan logika bisnis tidak bergantung langsung
// pada detail implementasi observability — aplikasi hanya perlu
// memanggil observability.Init() tanpa tahu teknologi di baliknya.
func Init() {
	log.Println("Initializing observability...")

	// Menginisialisasi komponen tracing (OpenTelemetry → Alloy → Tempo).
	// Tracing mencakup HTTP handler, service layer, dan DB query level.
	// Jika inisialisasi gagal, aplikasi tetap berjalan
	// tanpa menghentikan proses utama.
	if err := initTracing(); err != nil {
		log.Printf("Tracing init failed: %v", err)
	} else {
		log.Println("Tracing initialized (OTLP HTTP → Alloy → Tempo)")
	}

	// Menginisialisasi komponen profiling (Grafana Pyroscope).
	// Profiling bersifat opsional dan tidak mempengaruhi
	// fungsionalitas utama aplikasi.
	if err := initProfiling(); err != nil {
		log.Printf("Profiling init failed: %v", err)
	} else {
		log.Println("Profiling initialized (Pyroscope)")
	}

	// Menginisialisasi metrics tambahan.
	// Sebagian besar metrics didaftarkan langsung pada layer handler
	// menggunakan prometheus/client_golang.
	initMetrics()

	log.Println("Observability initialized successfully")
}
