import { ref } from 'vue'

const BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

export function useApi() {
  const loading = ref(false)
  const error = ref(null)

  async function request(path, options = {}) {
    loading.value = true
    error.value = null
    try {
      const response = await fetch(`${BASE_URL}${path}`, {
        headers: { 'Content-Type': 'application/json', ...options.headers },
        ...options,
      })
      if (!response.ok) {
        const body = await response.json().catch(() => ({}))
        throw new Error(body.error || `HTTP ${response.status}`)
      }
      return await response.json()
    } catch (err) {
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  const getProducts = (category) => request(`/api/products${category ? `?category=${category}` : ''}`)
  const getProduct = (id) => request(`/api/products/${id}`)
  const createOrder = (payload) => request('/api/orders', { method: 'POST', body: JSON.stringify(payload) })
  const getOrder = (id) => request(`/api/orders/${id}`)
  const getUserOrders = (userId) => request(`/api/users/${userId}/orders`)

  return { loading, error, getProducts, getProduct, createOrder, getOrder, getUserOrders }
}
