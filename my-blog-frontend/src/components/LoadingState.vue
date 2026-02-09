<script setup>
defineProps({
  type: {
    type: String,
    default: 'card'
  }
})
</script>

<template>
  <div class="loading-state">
    <el-skeleton :animated="true">
      <template v-if="type === 'card'" #template>
        <div class="skeleton-card">
          <el-skeleton-item variant="image" style="width: 100%; height: 200px" />
          <div class="skeleton-content">
            <el-skeleton-item variant="h3" style="width: 70%" />
            <el-skeleton-item variant="text" style="margin-top: 10px" />
            <el-skeleton-item variant="text" style="width: 60%" />
          </div>
        </div>
      </template>
      <template v-else-if="type === 'article'" #template>
        <div class="skeleton-article">
          <el-skeleton-item variant="h1" style="width: 80%; margin-bottom: 20px;" />
          <el-skeleton-item variant="text" style="margin-bottom: 20px;" />
          <el-skeleton-item variant="p" style="width: 100%; height: 400px;" />
        </div>
      </template>
      <template v-else-if="type === 'list'" #template>
        <div class="skeleton-list">
          <div v-for="i in 6" :key="i" class="skeleton-item">
            <el-skeleton-item variant="image" style="width: 100%; height: 200px" />
            <div class="skeleton-content">
              <el-skeleton-item variant="h3" style="width: 70%" />
              <el-skeleton-item variant="text" style="margin-top: 10px" />
            </div>
          </div>
        </div>
      </template>
    </el-skeleton>
  </div>
</template>

<style scoped>
.skeleton-card {
  background: #f8fafc;
  border-radius: 12px;
  overflow: hidden;
}

.skeleton-content {
  padding: 20px;
}

.skeleton-article {
  max-width: 800px;
  margin: 40px auto;
}

.skeleton-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
}

.skeleton-item {
  margin-bottom: 20px;
}
</style>
