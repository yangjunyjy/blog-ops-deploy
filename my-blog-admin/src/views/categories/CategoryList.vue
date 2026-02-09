<template>
  <div class="category-list">
    <el-card shadow="never">
      <!-- 搜索栏 -->
      <el-form :model="searchForm" class="search-form" @submit.prevent>
        <el-form-item label="名称">
          <el-input v-model="searchForm.name" placeholder="请输入名称" clearable style="width: 200px" />
        </el-form-item>
        <el-form-item class="search-actions">
          <HaloButton type="primary" size="small" content="搜索" :icon="Search" @click.prevent="handleSearch" />
          <HaloButton type="default" size="small" content="重置" :icon="Refresh" @click.prevent="handleReset" />
        </el-form-item>
      </el-form>

      <!-- 操作栏 -->
      <div class="action-bar">
        <HaloButton type="primary" size="medium" content="新增分类" :icon="Plus" @click="handleCreate" />
      </div>

      <!-- 表格 -->
      <el-table :data="categoryList" v-loading="loading" border stripe>
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="icon" label="图标" width="80" align="center">
          <template #default="{ row }">
            <span class="icon-text">{{ row.icon || '' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="name" label="名称" min-width="150" />
        <el-table-column prop="articleCount" label="文章数" width="100" align="center" />
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button type="danger" link size="small" @click="handleDelete(row.id)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <el-pagination v-model:current-page="pagination.page" v-model:page-size="pagination.pageSize"
        :page-sizes="[10, 20, 50, 100]" :total="pagination.total" layout="total, sizes, prev, pager, next, jumper"
        @size-change="loadList" @current-change="loadList" />
    </el-card>

    <!-- 编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="分类名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入分类名称" />
        </el-form-item>
        <el-form-item label="分类标识" prop="slug">
          <el-input v-model="form.slug" placeholder="请输入分类标识（英文）" />
        </el-form-item>
        <el-form-item label="分类图标" prop="icon">
          <el-input v-model="form.icon" placeholder="请输入图标（emoji）" />
        </el-form-item>
        <el-form-item label="排序" prop="sort_order">
          <el-input-number v-model="form.sort_order" :min="0" placeholder="请输入排序" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
defineOptions({
  name: 'CategoryList'
})

import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search, Refresh } from '@element-plus/icons-vue'
import {
  getCategoryList,
  createCategory,
  updateCategory,
  deleteCategory
} from '@/api/category'
import HaloButton from '@/layout/components/HaloButton.vue'

const loading = ref(false)
const dialogVisible = ref(false)
const submitting = ref(false)
const formRef = ref(null)
const categoryList = ref([])

const searchForm = reactive({
  name: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const dialogTitle = ref('新增分类')
const form = reactive({
  id: null,
  name: '',
  slug: '',
  icon: '💻',
  sort_order: 0
})

const rules = {
  name: [{ required: true, message: '请输入分类名称', trigger: 'blur' }],
  slug: [{ required: true, message: '请输入分类标识', trigger: 'blur' }]
}

const loadList = async () => {
  loading.value = true
  try {
    const res = await getCategoryList({
      page: pagination.page,
      pageSize: pagination.pageSize
    })
    categoryList.value = res.data.items || []
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
  searchForm.name = ''
  handleSearch()
}

const handleCreate = () => {
  dialogTitle.value = '新增分类'
  Object.assign(form, {
    id: null,
    name: '',
    slug: '',
    icon: '💻',
    sort_order: 0
  })
  dialogVisible.value = true
}

const handleEdit = (row) => {
  dialogTitle.value = '编辑分类'
  Object.assign(form, {
    id: row.id,
    name: row.name,
    slug: row.slug,
    icon: row.icon || '💻',
    sort_order: row.sort_order || 0
  })
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitting.value = true
    try {
      if (form.id) {
        await updateCategory(form.id, form)
        ElMessage.success('更新成功')
      } else {
        await createCategory(form)
        ElMessage.success('创建成功')
      }
      dialogVisible.value = false
      loadList()
    } catch (error) {
      ElMessage.error('操作失败')
    } finally {
      submitting.value = false
    }
  })
}

const handleDelete = async (id) => {
  try {
    await ElMessageBox.confirm('确定删除该分类吗？', '提示', { type: 'warning' })
    await deleteCategory(id)
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
.category-list {
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

  .icon-text {
    font-size: 24px;
  }
}
</style>
