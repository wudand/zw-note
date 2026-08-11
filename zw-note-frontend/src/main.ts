import { createApp } from 'vue'
import { createPinia } from 'pinia'
import './style.css'
import App from './App.vue'
import router from './router/index.ts'
import Components from '@/components'

// 已对接真实后端，默认关闭 mock；需要本地 mock 时设 VITE_USE_MOCK=true
if (import.meta.env.VITE_USE_MOCK === 'true') {
  await import('@/mock')
}

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(Components)
app.mount('#app')
