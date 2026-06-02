import { initializeFaro, getWebInstrumentations } from '@grafana/faro-web-sdk'
import { TracingInstrumentation } from '@grafana/faro-web-tracing'

let faroInstance = null

export function initFaro() {
  faroInstance = initializeFaro({
    url: import.meta.env.VITE_FARO_COLLECTOR_URL || 'http://alloy.monitoring.svc.cluster.local:12347/collect',
    app: {
      name: import.meta.env.VITE_APP_NAME || 'ecommerce-frontend',
      version: import.meta.env.VITE_APP_VERSION || '1.0.0',
      environment: import.meta.env.VITE_APP_ENV || 'production',
    },
    instrumentations: [
      ...getWebInstrumentations({ captureConsole: true }),
      new TracingInstrumentation({
        instrumentationOptions: {
          propagateTraceHeaderCorsUrls: [
            new RegExp(import.meta.env.VITE_API_URL || 'http://localhost:8080'),
          ],
        },
      }),
    ],
  })
  return faroInstance
}

export function getFaro() {
  return faroInstance
}
