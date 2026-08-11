import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  server: {
    port: 5555,
    host: '0.0.0.0',
    open: false,
    proxy: {
      '/api': {
        target: 'http://127.0.0.1:8004',
        changeOrigin: true,
      },
      // 图片等本地上传文件，由后端 /uploads 静态路由提供
      '/uploads': {
        target: 'http://127.0.0.1:8004',
        changeOrigin: true,
      },
    },
  }
})
