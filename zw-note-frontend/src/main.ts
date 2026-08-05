import '@/mock'
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import './style.css'
import App from './App.vue'
import router from './router/index.ts'
import Components from '@/components'

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(Components)
app.mount('#app')
