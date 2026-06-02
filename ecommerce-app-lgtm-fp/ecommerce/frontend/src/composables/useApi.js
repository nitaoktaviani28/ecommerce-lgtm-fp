/**
 * composables/useApi.js
 *
 * Composable untuk HTTP request ke BE.
 * Faro secara otomatis menginjeksikan W3C TraceContext headers
 * ke setiap fetch request, sehingga trace FE dan BE terhubung di Tempo.
 *
 * Tidak perlu ada kode tambahan untuk distributed tracing —
 * TracingInstrumentation di faro.js sudah menanganinya otomatis.
 */

import { ref } from 'vue'

const BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

export function useApi() {
  const loading = ref(false)
  const error = ref(null)

  /**
   * request adalah wrapper fetch yang otomatis:
   * 1. Set Content-Type header
   * 2. Handle error response
   * 3. Faro TracingInstrumentation inject traceparent header otomatis
   *    → trace FE terhubung ke trace BE di Tempo
   */
  async function request(path, options = {}) {
    loading.value = true
    error.value = null

    try {
      const response = await fetch(`${BASE_URL}${path}`, {
        headers: {
          'Content-Type': 'application/json',
          ...options.headers,
        },
        ...options,
      })

      if (!response.ok) {
        const body = await response.json().catch(() => ({}))
        throw new Error(body.error || `HTTP ${response.status}`)
      }

      return await response.json()
    } catch (err) {
      error.value = err.message
      // Faro otomatis menangkap error ini dan mengirimnya ke Loki
      throw err
    } finally {
      loading.value = false
    }
  }

  // ─── Product APIs ──────────────────────────────────────────────────────────
  const getProducts = (category) =>
    request(`/api/products${category ? `?category=${category}` : ''}`)

  const getProduct = (id) =>
    request(`/api/products/${id}`)

  // ─── Order APIs ────────────────────────────────────────────────────────────
  const createOrder = (payload) =>
    request('/api/orders', {
      method: 'POST',
      body: JSON.stringify(payload),
    })

  const getOrder = (id) =>
    request(`/api/orders/${id}`)

  const getUserOrders = (userId) =>
    request(`/api/users/${userId}/orders`)

  return {
    loading,
    error,
    getProducts,
    getProduct,
    createOrder,
    getOrder,
    getUserOrders,
  }
}
