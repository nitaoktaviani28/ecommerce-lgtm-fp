import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { getFaro } from '@/observability/faro'

export const useCartStore = defineStore('cart', () => {
  const items = ref([])
  const totalItems = computed(() => items.value.reduce((sum, i) => sum + i.quantity, 0))
  const totalPrice = computed(() => items.value.reduce((sum, i) => sum + i.product.price * i.quantity, 0))

  function addItem(product, quantity = 1) {
    const existing = items.value.find(i => i.product.id === product.id)
    if (existing) {
      existing.quantity += quantity
    } else {
      items.value.push({ product, quantity })
    }
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
