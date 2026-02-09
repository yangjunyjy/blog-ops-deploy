<template>
  <div class="comment-list">
    <el-card shadow="never">
      <!-- 搜索栏 -->
      <el-form :model="searchForm" class="search-form" @submit.prevent>
        <el-form-item label="文章ID">
          <el-input v-model="searchForm.articleId" placeholder="请输入文章ID" clearable style="width: 150px" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="请选择" clearable style="width: 120px" :teleported="false">
            <el-option label="全部" value="" />
            <el-option label="待审核" :value="0" />
            <el-option label="已通过" :value="1" />
            <el-option label="已拒绝" :value="2" />
          </el-select>
        </el-form-item>
        <el-form-item class="search-actions">
          <HaloButton type="primary" size="small" content="搜索" :icon="Search" @click.prevent="handleSearch" />
          <HaloButton type="default" size="small" content="重置" :icon="Refresh" @click.prevent="handleReset" />
        </el-form-item>
      </el-form>

      <!-- 表格 -->
      <el-table :data="commentList" v-loading="loading" border stripe>
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column label="用户" width="150">
          <template #default="{ row }">
            <div class="user-info">
              <el-avatar :size="32" :src="row.user?.avatar" />
              <span class="username">{{ row.user?.nickname }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="文章" min-width="200" show-overflow-tooltip>
          <template #default="{ row }">
            {{ row.article?.title }}
          </template>
        </el-table-column>
        <el-table-column prop="content" label="评论内容" min-width="250" show-overflow-tooltip />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="创建时间" width="160">
          <template #default="{ row }">
            {{ formatDate(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button v-if="row.status === 0" type="success" link size="small" @click="handleApprove(row.id)">
              通过
            </el-button>
            <el-button v-if="row.status === 0" type="warning" link size="small" @click="handleReject(row.id)">
              拒绝
            </el-button>
            <el-button type="danger" link size="small" @click="handleDelete(row.id)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <el-pagination v-model:current-page="pagination.page" v-model:page-size="pagination.size"
        :page-sizes="[10, 20, 50, 100]" :total="pagination.total" layout="total, sizes, prev, pager, next, jumper"
        @size-change="loadList" @current-change="loadList" style="margin-top: 20px; justify-content: flex-end" />
    </el-card>
  </div>
</template>

<script setup>
defineOptions({
  name: 'CommentList'
})

import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Refresh } from '@element-plus/icons-vue'
import { getCommentList, deleteComment } from '@/api'
import HaloButton from '@/layout/components/HaloButton.vue'

const loading = ref(false)
const commentList = ref([])

const searchForm = reactive({
  articleId: '',
  status: 1
})

const pagination = reactive({
  page: 1,
  size: 10,
  total: 0
})

const formatDate = (date) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}

const getStatusText = (status) => {
  const map = { 0: '待审核', 1: '已通过', 2: '已拒绝' }
  return map[status] || '未知'
}

const getStatusType = (status) => {
  const map = { 0: 'warning', 1: 'success', 2: 'info' }
  return map[status] || 'info'
}

const loadList = async () => {
  loading.value = true
  try {
    const res = await getCommentList({
      ...searchForm,
      page: pagination.page,
      pageSize: pagination.size
    })
    commentList.value = res.data.items || []
    pagination.total = res.data.total || 0
  } catch (error) {
    ElMessage.error('加载失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  loadList()
}

const handleReset = () => {
  searchForm.articleId = ''
  searchForm.status = ''
  handleSearch()
}

const handleApprove = async (id) => {
  try {
    await mockApi.approveComment(id)
    ElMessage.success('审核通过')
    loadList()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

const handleReject = async (id) => {
  try {
    await mockApi.rejectComment(id)
    ElMessage.success('审核拒绝')
    loadList()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

const handleDelete = async (id) => {
  try {
    await ElMessageBox.confirm('确定删除该评论吗？', '提示', { type: 'warning' })
    await deleteComment(id)
    ElMessage.success('删除成功')
    loadList()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

onMounted(() => {
  loadList()
})
</script>

<style scoped lang="scss">
.comment-list {
  height: 100%;

  :deep(.el-card) {
    height: 100%;
    display: flex;
    flex-direction: column;

    .el-card__body {
      height: 100%;
      display: flex;
      flex-direction: column;
      padding: 16px;
    }
  }

  .search-form {
    margin-bottom: 12px;
    display: flex;
    align-items: flex-start;
    flex-wrap: wrap;
    gap: 12px;
    flex-shrink: 0;

    :deep(.el-form-item) {
      margin-bottom: 0;
    }

    .search-actions {
      margin-left: auto;

      :deep(.el-form-item__content) {
        display: flex;
      }

      .HaloButton {
        margin-left: 8px;
      }
    }
  }

  :deep(.el-table) {
    flex: 1;
    overflow: auto;
  }

  :deep(.el-table__body-wrapper) {
    overflow: auto;
  }

  :deep(.el-pagination) {
    margin-top: 12px;
    justify-content: flex-end;
    flex-shrink: 0;
  }

  .user-info {
    display: flex;
    align-items: center;
    gap: 10px;

    .username {
      font-size: 14px;
      color: #303133;
    }
  }
}
</style>
