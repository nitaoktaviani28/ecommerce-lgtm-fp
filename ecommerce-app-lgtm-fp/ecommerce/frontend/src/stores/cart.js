/**
 * stores/cart.js — Pinia store untuk shopping cart.
 */

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { getFaro } from '@/observability/faro'

export const useCartStore = defineStore('cart', () => {
  const items = ref([]) // [{ product, quantity }]

  const totalItems = computed(() =>
    items.value.reduce((sum, item) => sum + item.quantity, 0)
  )

  const totalPrice = computed(() =>
    items.value.reduce((sum, item) => sum + item.product.price * item.quantity, 0)
  )

  function addItem(product, quantity = 1) {
    const existing = items.value.find(i => i.product.id === product.id)
    if (existing) {
      existing.quantity += quantity
    } else {
      items.value.push({ product, quantity })
    }

    // Kirim custom event ke Faro → Loki
    // Event ini bisa dipakai untuk analisis funnel di Grafana
    getFaro()?.api.pushEvent('cart_item_added', {
      product_id: String(product.id),
      product_name: product.name,
      quantity: String(quantity),
    })
  }

  function removeItem(productId) {
    items.value = items.value.filter(i => i.product.id !== productId)
  }

  function clearCart() {
    items.value = []
  }

  return { items, totalItems, totalPrice, addItem, removeItem, clearCart }
})
