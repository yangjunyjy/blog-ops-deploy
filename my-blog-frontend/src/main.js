import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import VueLazyload from 'vue3-lazyload'
import App from './App.vue'
import router from './router'
import { formatDate, formatRelativeTime } from './utils/format'
import './styles/global.css'

const app = createApp(App)

// 注册全局过滤器
app.config.globalProperties.$filters = {
  formatDate,
  formatRelativeTime
}

// 注册所有图标
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

app.use(createPinia())
app.use(router)
app.use(ElementPlus)
app.use(VueLazyload, {
  loading: '/loading.gif',
  error: '/error.gif'
})

app.mount('#app')
