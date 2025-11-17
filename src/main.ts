import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import router from './router'

import Components from '@/components';

createApp(App).use(router).use(Components).mount('#app')
