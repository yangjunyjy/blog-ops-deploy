<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import Header from '@/components/Header.vue'
import Footer from '@/components/Footer.vue'

const route = useRoute()

// 判断是否隐藏布局组件（导航栏和页脚）
const hideLayout = computed(() => route.meta.hideLayout || false)
</script>

<template>
  <div id="app">
    <Header v-if="!hideLayout" />
    <main class="main-content" :class="{ 'no-layout': hideLayout }">
      <router-view v-slot="{ Component }">
        <keep-alive>
          <component :is="Component" v-if="$route.meta.keepAlive" />
        </keep-alive>
        <component :is="Component" v-if="!$route.meta.keepAlive" />
      </router-view>
    </main>
    <Footer v-if="!hideLayout" />
  </div>
</template>

<style>
#app {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

.main-content {
  flex: 1;
}

.main-content.no-layout {
  flex: 1;
  min-height: 100vh;
}
</style>
