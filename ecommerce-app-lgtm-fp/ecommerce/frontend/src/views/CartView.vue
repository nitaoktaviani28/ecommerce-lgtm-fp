<template>
  <div class="cart">
    <h2>🛒 Keranjang Belanja</h2>
    <div v-if="cartStore.items.length === 0" class="empty">
      <p>Keranjang kosong</p>
      <button class="btn-primary" @click="$router.push('/')">Belanja Sekarang</button>
    </div>
    <div v-else>
      <div v-for="item in cartStore.items" :key="item.product.id" class="cart-item">
        <div class="item-info">
          <h4>{{ item.product.name }}</h4>
          <span class="item-category">{{ item.product.category }}</span>
        </div>
        <div class="item-controls">
          <span class="item-price">Rp {{ formatPrice(item.product.price * item.quantity) }}</span>
          <span class="item-qty">× {{ item.quantity }}</span>
          <button class="btn-remove" @click="cartStore.removeItem(item.product.id)">✕</button>
        </div>
      </div>
      <div class="cart-summary">
        <div class="total">
          <span>Total</span>
          <span class="total-price">Rp {{ formatPrice(cartStore.totalPrice) }}</span>
        </div>
        <div class="user-id-input">
          <label>User ID:</label>
          <input v-model="userId" type="number" placeholder="Contoh: 1" />
        </div>
        <button class="btn-checkout" :disabled="loading || !userId" @click="checkout">
          {{ loading ? 'Memproses...' : 'Checkout Sekarang' }}
        </button>
        <div v-if="error" class="error">{{ error }}</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useCartStore } from '@/stores/cart'
import { useApi } from '@/composables/useApi'
import { getFaro } from '@/observability/faro'

const router = useRouter()
const cartStore = useCartStore()
const { createOrder, loading, error } = useApi()
const userId = ref(1)

async function checkout() {
  getFaro()?.api.pushEvent('checkout_initiated', {
    total_price: String(cartStore.totalPrice),
    item_count: String(cartStore.totalItems),
  })
  try {
    const order = await createOrder({
      user_id: Number(userId.value),
      items: cartStore.items.map(i => ({ product_id: i.product.id, quantity: i.quantity })),
    })
    getFaro()?.api.pushEvent('checkout_completed', { order_id: String(order.id) })
    cartStore.clearCart()
    router.push(`/orders/${order.id}`)
  } catch {}
}

function formatPrice(price) {
  return new Intl.NumberFormat('id-ID').format(price)
}
</script>

<style scoped>
.cart { max-width: 600px; margin: 40px auto; padding: 20px; }
.empty { text-align: center; padding: 60px 0; }
.cart-item { display: flex; justify-content: space-between; align-items: center; padding: 16px; border: 1px solid #e5e7eb; border-radius: 10px; margin-bottom: 12px; background: white; }
.item-info h4 { margin: 0 0 4px; }
.item-category { font-size: 12px; color: #9ca3af; }
.item-controls { display: flex; align-items: center; gap: 12px; }
.item-price { font-weight: 700; color: #4f46e5; }
.item-qty { color: #6b7280; }
.btn-remove { background: none; border: 1px solid #e5e7eb; border-radius: 6px; padding: 4px 8px; cursor: pointer; color: #ef4444; }
.cart-summary { margin-top: 24px; padding: 20px; background: #f9fafb; border-radius: 12px; }
.total { display: flex; justify-content: space-between; font-size: 1.2rem; margin-bottom: 16px; }
.total-price { font-weight: 700; color: #4f46e5; }
.user-id-input { margin-bottom: 16px; }
.user-id-input label { display: block; margin-bottom: 6px; font-size: 14px; color: #6b7280; }
.user-id-input input { width: 100%; padding: 10px; border: 1px solid #e5e7eb; border-radius: 8px; font-size: 16px; }
.btn-checkout { width: 100%; padding: 14px; background: #10b981; color: white; border: none; border-radius: 10px; font-size: 1rem; font-weight: 600; cursor: pointer; }
.btn-checkout:disabled { opacity: 0.6; cursor: not-allowed; }
.btn-primary { padding: 12px 24px; background: #4f46e5; color: white; border: none; border-radius: 8px; cursor: pointer; margin-top: 16px; }
.error { color: #ef4444; margin-top: 12px; text-align: center; }
</style>
