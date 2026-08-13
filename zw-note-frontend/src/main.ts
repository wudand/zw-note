import { createApp } from 'vue'
import { createPinia } from 'pinia'
import 'highlight.js/styles/github.css'
import App from './App.vue'
import router from './router/index.ts'
import Components from '@/components'
// 必须最后导入：Element Plus 的默认 CSS 变量（如 --el-color-primary-light-3）
// 是在上面 Components 里引入的，只有 style.css 在其之后加载，
// 才能保证同优先级、未加 !important 的变量覆盖也能生效（例如按钮 hover/active 配色）。
import './style.css'

// 已对接真实后端，默认关闭 mock；需要本地 mock 时设 VITE_USE_MOCK=true
if (import.meta.env.VITE_USE_MOCK === 'true') {
  await import('@/mock')
}

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(Components)
app.mount('#app')
