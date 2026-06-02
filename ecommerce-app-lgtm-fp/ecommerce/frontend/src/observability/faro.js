/**
 * observability/faro.js
 *
 * Modul ini menginisialisasi Grafana Faro untuk Frontend Observability (RUM).
 * Sama seperti pola di BE (observability.Init()), modul ini menjadi
 * satu-satunya titik masuk observability di FE.
 *
 * Yang dikumpulkan Faro secara otomatis:
 * - JavaScript errors & exceptions
 * - Web Vitals (LCP, FID, CLS, TTFB)
 * - User sessions & interactions
 * - Console logs (warn, error)
 * - Distributed traces (terhubung ke trace BE via W3C TraceContext headers)
 *
 * Komponen Vue lain tidak perlu tahu detail implementasi ini.
 */

import {
  initializeFaro,
  getWebInstrumentations,
} from '@grafana/faro-web-sdk'

import { TracingInstrumentation } from '@grafana/faro-web-tracing'

let faroInstance = null

/**
 * initFaro menginisialisasi Grafana Faro.
 * Dipanggil sekali di main.js sebelum app.mount().
 *
 * Faro akan otomatis:
 * 1. Menangkap semua JS errors → dikirim ke Loki via Alloy
 * 2. Mengukur Web Vitals → dikirim ke Mimir via Alloy
 * 3. Membuat traces dari browser → dikirim ke Tempo via Alloy
 *    dan otomatis terhubung ke trace BE melalui W3C TraceContext headers
 */
export function initFaro() {
  const collectorUrl = import.meta.env.VITE_FARO_COLLECTOR_URL ||
    'http://alloy.monitoring.svc.cluster.local:12347/collect'

  faroInstance = initializeFaro({
    // URL endpoint Faro collector (Alloy dengan Faro receiver)
    url: collectorUrl,

    // Identitas aplikasi di Grafana
    app: {
      name: import.meta.env.VITE_APP_NAME || 'ecommerce-frontend',
      version: import.meta.env.VITE_APP_VERSION || '1.0.0',
      environment: import.meta.env.VITE_APP_ENV || 'production',
    },

    // Instrumentasi bawaan yang diaktifkan
    instrumentations: [
      // Web Vitals, errors, console logs, sessions, performance
      ...getWebInstrumentations({
        captureConsole: true,           // tangkap console.warn dan console.error
        captureConsoleDisabledLevels: ['log', 'debug', 'info'], // skip noise
      }),

      // Distributed tracing dari browser → BE
      // Faro akan inject W3C TraceContext headers ke semua fetch/XHR request
      // sehingga trace FE dan BE terhubung dalam satu waterfall di Tempo
      new TracingInstrumentation({
        instrumentationOptions: {
          // Propagate trace context ke BE API
          propagateTraceHeaderCorsUrls: [
            new RegExp(import.meta.env.VITE_API_URL || 'http://localhost:8080'),
          ],
        },
      }),
    ],
  })

  console.info('[Faro] Frontend observability initialized')
  return faroInstance
}

/**
 * getFaro mengembalikan instance Faro yang sudah diinisialisasi.
 * Gunakan ini untuk manual instrumentation jika diperlukan.
 *
 * Contoh penggunaan di komponen Vue:
 *   import { getFaro } from '@/observability/faro'
 *   getFaro()?.api.pushEvent('checkout_clicked', { product_id: '123' })
 */
export function getFaro() {
  return faroInstance
}
