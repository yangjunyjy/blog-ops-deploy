import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  server: {
    port: 3000,
    open: true,
    proxy: {
      '/api': {
        target: 'http://localhost:8081',
        changeOrigin: true
      },
      // WebSocket 代理配置（需要重写路径为 /api/v1/rbac）
      '/rbac': {
        target: 'http://localhost:8081',
        changeOrigin: true,
        ws: true,
        rewrite: (path) => {
          return '/api/v1' + path
        }
      }
    }
  }
})
