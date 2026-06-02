package observability

import (
	"context"
	"database/sql"

	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	otelsemconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

// initTracing menginisialisasi sistem tracing menggunakan OpenTelemetry.
// Fungsi ini menyiapkan exporter, identitas service, dan tracer provider,
// kemudian mendaftarkannya secara global agar dapat digunakan
// oleh seluruh bagian aplikasi termasuk layer database.
func initTracing() error {
	// Membuat exporter OTLP HTTP untuk mengirim trace ke Alloy → Tempo.
	// Endpoint diambil dari environment variable agar fleksibel
	// untuk berbagai environment (local, staging, production).
	exporter, err := otlptracehttp.New(
		context.Background(),
		otlptracehttp.WithEndpoint(
			GetEnv(
				"OTEL_EXPORTER_OTLP_ENDPOINT",
				"http://alloy.monitoring.svc.cluster.local:4318",
			),
		),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return err
	}

	// Mendefinisikan resource OpenTelemetry sebagai identitas service.
	// service.name digunakan oleh Tempo dan Grafana untuk mengelompokkan trace.
	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			otelsemconv.ServiceName(
				GetEnv("OTEL_SERVICE_NAME", "ecommerce-backend"),
			),
			otelsemconv.DeploymentEnvironment(
				GetEnv("APP_ENV", "production"),
			),
		),
	)
	if err != nil {
		return err
	}

	// Membuat TracerProvider sebagai mesin utama tracing.
	// Batcher digunakan agar pengiriman trace lebih efisien.
	// AlwaysSample digunakan agar semua request di-trace — cocok untuk learning.
	// Untuk production heavy traffic, ganti ke sdktrace.TraceIDRatioBased(0.1)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	// Mendaftarkan TracerProvider secara global.
	// Setelah ini, seluruh tracer di aplikasi akan menggunakan provider ini.
	otel.SetTracerProvider(tp)

	return nil
}

// OpenDBWithTracing membuka koneksi database PostgreSQL dengan instrumentasi
// OpenTelemetry otomatis menggunakan otelsql.
//
// Setiap query SQL yang dieksekusi melalui koneksi ini akan otomatis
// menghasilkan span tracing dengan informasi:
//   - Nama operasi SQL (SELECT, INSERT, UPDATE, DELETE)
//   - Query string yang dieksekusi
//   - Nama tabel yang diakses
//   - Status sukses atau error
//
// Dengan ini, trace dari HTTP handler → service → repository → DB query
// semuanya terhubung dalam satu distributed trace di Tempo.
func OpenDBWithTracing(dsn string) (*sql.DB, error) {
	db, err := otelsql.Open(
		"postgres",
		dsn,
		// Mendaftarkan atribut database ke dalam span tracing.
		// Ini membuat Tempo menampilkan info koneksi DB pada setiap span.
		otelsql.WithAttributes(
			semconv.DBSystemPostgreSQL,
		),
		// Mengaktifkan perekaman nama operasi SQL (SELECT, INSERT, dll).
		otelsql.WithDBName(
			GetEnv("DB_NAME", "ecommerce"),
		),
	)
	if err != nil {
		return nil, err
	}

	// Mendaftarkan metrics koneksi DB ke Prometheus.
	// Ini menghasilkan metrics seperti:
	// - db_sql_connections_open
	// - db_sql_connections_idle
	// - db_sql_connections_wait_duration
	if err := otelsql.RegisterDBStatsMetrics(db,
		otelsql.WithAttributes(semconv.DBSystemPostgreSQL),
	); err != nil {
		return nil, err
	}

	return db, nil
}
