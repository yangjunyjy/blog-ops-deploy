<template>
  <div class="article-list">
    <el-card shadow="never">
      <!-- 搜索栏 -->
      <el-form :model="searchForm" class="search-form" @submit.prevent>
        <el-form-item label="标题">
          <el-input v-model="searchForm.title" placeholder="请输入标题" clearable style="width: 200px" />
        </el-form-item>
        <el-form-item label="分类">
          <el-select v-model="searchForm.categoryId" placeholder="请选择" clearable style="width: 150px">
            <el-option v-for="cat in categories" :key="cat.id" :label="cat.name" :value="cat.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable style="width: 120px" :teleported="false">
            <el-option label="已发布" value="1" />
            <el-option label="草稿" value="0" />
          </el-select>
        </el-form-item>
        <el-form-item class="search-actions">
          <HaloButton type="primary" size="small" content="搜索" :icon="Search" @click.prevent="handleSearch" />
          <HaloButton type="default" size="small" content="重置" :icon="Refresh" @click.prevent="handleReset" />
        </el-form-item>
      </el-form>

      <!-- 操作栏 -->
      <div class="action-bar">
        <HaloButton type="primary" size="medium" content="写文章" :icon="Plus" @click="handleCreate" />
        <HaloButton type="warning" size="medium" :content="`导出文章 (${selectedIds.length})`" :icon="Download"
          @click="handleExport" :disabled="selectedIds.length === 0" />
      </div>

      <!-- 表格 -->
      <el-table :data="articleList" v-loading="loading" border stripe @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" />
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column label="封面" width="100">
          <template #default="{ row }">
            <el-image v-if="row.cover" :src="row.cover" fit="cover"
              style="width: 80px; height: 50px; border-radius: 4px" />
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="title" label="标题" min-width="200" show-overflow-tooltip />
        <el-table-column label="分类" width="120">
          <template #default="{ row }">
            <el-tag v-if="row.category" type="success" size="small">{{ row.category.name }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="标签" width="180">
          <template #default="{ row }">
            <el-tag v-for="(tag, i) in (row.tags || []).slice(0, 3)" :key="i" type="info" size="small"
              style="margin-right: 4px">
              {{ tag.name }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="views" label="浏览" width="80" align="center" />
        <el-table-column label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
              {{ row.status === 1 ? '已发布' : '草稿' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="创建时间" width="160">
          <template #default="{ row }">
            {{ formatDate(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button type="success" link size="small" @click="handleDownload(row)">下载</el-button>
            <el-button type="danger" link size="small" @click="handleDelete(row.id)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <el-pagination v-model:current-page="pagination.page" v-model:page-size="pagination.size"
        :page-sizes="[10, 20, 50, 100]" :total="pagination.total" layout="total, sizes, prev, pager, next, jumper"
        @size-change="loadList" @current-change="loadList" />
    </el-card>
  </div>
</template>

<script setup>
defineOptions({
  name: 'ArticleList'
})

import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Download, Search, Refresh } from '@element-plus/icons-vue'
import { mockApi } from '../../mock'
import { useRbacStore } from '../../store/rbac'
import { getArticleList, deleteArticle } from '@/api'
import HaloButton from '@/layout/components/HaloButton.vue'

const router = useRouter()

const loading = ref(false)
const articleList = ref([])
const categories = ref([])
const selectedIds = ref([])
const selectedArticles = ref([])
const rbacStore = useRbacStore()

const searchForm = reactive({
  title: '',
  categoryId: '',
  status: '1'
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

const loadList = async () => {
  loading.value = true
  try {
    const res = await getArticleList({
      ...searchForm,
      page: pagination.page,
      pageSize: pagination.size
    })
    articleList.value = res.data.items || []
    pagination.total = res.data.total || 0
  } catch (error) {
    ElMessage.error('加载失败')
  } finally {
    loading.value = false
  }
}

const loadCategories = async () => {
  try {
    const res = await mockApi.getCategoryList()
    categories.value = res.data.list || []
  } catch (error) {
    console.error(error)
  }
}

const handleSearch = () => {
  pagination.page = 1
  loadList()
}

const handleReset = () => {
  Object.assign(searchForm, {
    title: '',
    categoryId: '',
    status: ''
  })
  handleSearch()
}

const handleCreate = () => {
  if (rbacStore.hasPermission('content:article:create')) {
    router.push('/create')
  } else {
    ElMessage.warning('您没有创建文章的权限')
  }
}

const handleEdit = (row) => {
  router.push(`/editor/${row.id}`)
}

const handleDelete = async (id) => {
  try {
    await ElMessageBox.confirm('确定删除该文章吗？', '提示', { type: 'warning' })
    await deleteArticle(id)
    ElMessage.success('删除成功')
    loadList()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 表格选择变化
const handleSelectionChange = (selection) => {
  selectedArticles.value = selection
  selectedIds.value = selection.map(item => item.id)
}

// 导出选中的文章
const handleExport = async () => {
  if (selectedArticles.value.length === 0) {
    ElMessage.warning('请先选择要导出的文章')
    return
  }

  try {
    // 按批次导出（多个文件打包成 zip）
    const JSZip = (await import('jszip')).default
    const zip = new JSZip()

    selectedArticles.value.forEach((article, index) => {
      const fileName = `${article.title.replace(/[\\/:*?"<>|]/g, '_')}.md`
      const content = `# ${article.title}\n\n${article.content || ''}`
      zip.file(fileName, content)
    })

    // 生成 ZIP 文件
    const zipContent = await zip.generateAsync({ type: 'blob' })
    const url = URL.createObjectURL(zipContent)
    const link = document.createElement('a')
    link.href = url
    link.download = `articles_${new Date().getTime()}.zip`
    link.click()
    URL.revokeObjectURL(url)

    ElMessage.success(`成功导出 ${selectedArticles.value.length} 篇文章`)
  } catch (error) {
    ElMessage.error('导出失败')
    console.error(error)
  }
}

// 下载单篇文章
const handleDownload = async (article) => {
  try {
    const content = `# ${article.title}\n\n${article.content || ''}`
    const blob = new Blob([content], { type: 'text/markdown;charset=utf-8' })
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `${article.title.replace(/[\\/:*?"<>|]/g, '_')}.md`
    link.click()
    URL.revokeObjectURL(url)

    ElMessage.success('下载成功')
  } catch (error) {
    ElMessage.error('下载失败')
  }
}

onMounted(() => {
  loadCategories()
  loadList()
})
</script>

<style scoped lang="scss">
.article-list {
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

  .action-bar {
    margin-bottom: 12px;
    display: flex;
    gap: 12px;
    flex-shrink: 0;
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
}
</style>
