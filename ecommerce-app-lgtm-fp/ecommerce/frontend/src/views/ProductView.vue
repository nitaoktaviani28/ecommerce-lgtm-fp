<template>
  <div class="product-detail">
    <button class="btn-back" @click="$router.push('/')">← Kembali</button>
    <div v-if="loading" class="loading">Memuat produk...</div>
    <div v-else-if="error" class="error">{{ error }}</div>
    <div v-else-if="product" class="detail-card">
      <div class="product-category">{{ product.category }}</div>
      <h2>{{ product.name }}</h2>
      <p class="desc">{{ product.description }}</p>
      <div class="price-row">
        <span class="price">Rp {{ formatPrice(product.price) }}</span>
        <span :class="['stock', product.stock < 5 ? 'low' : '']">Stok: {{ product.stock }}</span>
      </div>
      <div class="qty-row">
        <label>Jumlah:</label>
        <input v-model.number="quantity" type="number" min="1" :max="product.stock" />
      </div>
      <button class="btn-add" :disabled="product.stock === 0" @click="addToCart">
        {{ product.stock === 0 ? 'Stok Habis' : '+ Tambah ke Keranjang' }}
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useApi } from '@/composables/useApi'
import { useCartStore } from '@/stores/cart'

const route = useRoute()
const router = useRouter()
const { getProduct, loading, error } = useApi()
const cartStore = useCartStore()
const product = ref(null)
const quantity = ref(1)

onMounted(async () => {
  try { product.value = await getProduct(route.params.id) } catch {}
})

function addToCart() {
  cartStore.addItem(product.value, quantity.value)
  router.push('/cart')
}

function formatPrice(price) {
  return new Intl.NumberFormat('id-ID').format(price)
}
</script>

<style scoped>
.product-detail { max-width: 600px; margin: 40px auto; padding: 20px; }
.btn-back { background: none; border: none; color: #4f46e5; cursor: pointer; font-size: 14px; margin-bottom: 20px; }
.detail-card { background: white; border: 1px solid #e5e7eb; border-radius: 16px; padding: 24px; }
.product-category { font-size: 12px; color: #6b7280; text-transform: uppercase; margin-bottom: 8px; }
.desc { color: #6b7280; line-height: 1.6; margin: 12px 0; }
.price-row { display: flex; justify-content: space-between; align-items: center; margin: 16px 0; }
.price { font-size: 1.5rem; font-weight: 700; color: #4f46e5; }
.stock { font-size: 14px; color: #10b981; }
.stock.low { color: #ef4444; }
.qty-row { display: flex; align-items: center; gap: 12px; margin-bottom: 16px; }
.qty-row input { width: 80px; padding: 8px; border: 1px solid #e5e7eb; border-radius: 8px; font-size: 16px; text-align: center; }
.btn-add { width: 100%; padding: 14px; background: #4f46e5; color: white; border: none; border-radius: 10px; font-size: 1rem; font-weight: 600; cursor: pointer; }
.btn-add:disabled { opacity: 0.5; cursor: not-allowed; }
.loading, .error { text-align: center; padding: 40px; color: #6b7280; }
</style>
