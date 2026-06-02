import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import App from './App.vue'
import { initFaro } from './observability/faro'

// 1. Init Faro sebelum app mount
initFaro()

// 2. Router
import HomeView from './views/HomeView.vue'
import ProductView from './views/ProductView.vue'
import CartView from './views/CartView.vue'
import OrderView from './views/OrderView.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', component: HomeView, name: 'home' },
    { path: '/products/:id', component: ProductView, name: 'product' },
    { path: '/cart', component: CartView, name: 'cart' },
    { path: '/orders/:id', component: OrderView, name: 'order' },
  ],
})

// 3. Mount
const app = createApp(App)
app.use(createPinia())
app.use(router)
app.mount('#app')
