/**
 * main.js — Entry point aplikasi Vue.js ecommerce.
 *
 * Observability (Faro) diinisialisasi di sini sebelum app mount,
 * sehingga seluruh error dan interaction sejak awal sudah ter-capture.
 * Komponen lain tidak perlu tahu detail implementasi Faro.
 */

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'

import App from './App.vue'
import { initFaro } from './observability/faro'

// ─── 1. Inisialisasi Faro (harus sebelum app.mount) ─────────────────────────
// Sama seperti observability.Init() di Go BE,
// ini adalah single entry point observability di FE.
initFaro()

// ─── 2. Router ───────────────────────────────────────────────────────────────
import HomeView from './views/HomeView.vue'
import ProductView from './views/ProductView.vue'
import CartView from './views/CartView.vue'
import OrderView from './views/OrderView.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/',               component: HomeView,    name: 'home' },
    { path: '/products/:id',   component: ProductView, name: 'product' },
    { path: '/cart',           component: CartView,    name: 'cart' },
    { path: '/orders/:id',     component: OrderView,   name: 'order' },
  ],
})

// ─── 3. Store ─────────────────────────────────────────────────────────────────
const pinia = createPinia()

// ─── 4. Mount App ─────────────────────────────────────────────────────────────
const app = createApp(App)
app.use(router)
app.use(pinia)
app.mount('#app')
