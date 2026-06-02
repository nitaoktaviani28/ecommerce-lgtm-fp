<template>
  <div class="order-detail">
    <h2>✅ Order Berhasil!</h2>
    <div v-if="loading" class="loading">Memuat order...</div>
    <div v-else-if="error" class="error">{{ error }}</div>
    <div v-else-if="order" class="order-card">
      <div class="order-header">
        <div>
          <span class="label">Order ID</span>
          <span class="value">#{{ order.id }}</span>
        </div>
        <span class="status">{{ order.status }}</span>
      </div>
      <div class="items">
        <div v-for="item in order.items" :key="item.id" class="item">
          <span>Produk #{{ item.product_id }}</span>
          <span>× {{ item.quantity }}</span>
          <span>Rp {{ formatPrice(item.price * item.quantity) }}</span>
        </div>
      </div>
      <div class="total">
        <span>Total</span>
        <span class="total-price">Rp {{ formatPrice(order.total_price) }}</span>
      </div>
      <button class="btn-primary" @click="$router.push('/')">Lanjut Belanja</button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useApi } from '@/composables/useApi'

const route = useRoute()
const { getOrder, loading, error } = useApi()
const order = ref(null)

onMounted(async () => {
  try { order.value = await getOrder(route.params.id) } catch {}
})

function formatPrice(price) {
  return new Intl.NumberFormat('id-ID').format(price)
}
</script>

<style scoped>
.order-detail { max-width: 600px; margin: 40px auto; padding: 20px; }
h2 { color: #10b981; margin-bottom: 24px; }
.order-card { background: white; border: 1px solid #e5e7eb; border-radius: 16px; padding: 24px; }
.order-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
.label { font-size: 12px; color: #6b7280; display: block; }
.value { font-size: 1.2rem; font-weight: 700; }
.status { padding: 6px 14px; border-radius: 20px; font-size: 12px; font-weight: 600; text-transform: uppercase; background: #d1fae5; color: #065f46; }
.items { border-top: 1px solid #f3f4f6; border-bottom: 1px solid #f3f4f6; padding: 16px 0; margin-bottom: 16px; }
.item { display: flex; justify-content: space-between; padding: 8px 0; font-size: 14px; }
.total { display: flex; justify-content: space-between; font-size: 1.1rem; margin-bottom: 20px; }
.total-price { font-weight: 700; color: #4f46e5; }
.btn-primary { width: 100%; padding: 14px; background: #4f46e5; color: white; border: none; border-radius: 10px; font-size: 1rem; cursor: pointer; }
.loading, .error { text-align: center; padding: 40px; color: #6b7280; }
</style>
