<template>
  <div class="tag-list">
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
        <HaloButton type="primary" size="medium" content="新增标签" :icon="Plus" @click="handleCreate" />
      </div>

      <!-- 表格 -->
      <el-table :data="tagList" v-loading="loading" border stripe>
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="name" label="名称" min-width="150" />
        <el-table-column prop="slug" label="URL后缀" min-width="150">
          <template #default="{ row }">
            <el-tag size="small">/tag/{{ row.slug }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
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
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px" @closed="handleDialogClosed">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="标签名称" prop="name" required>
          <el-input v-model="form.name" placeholder="请输入标签名称" maxlength="50" show-word-limit />
        </el-form-item>

        <el-form-item label="友好URL" prop="slug" required>
          <el-input v-model="form.slug" placeholder="请输入友好URL后缀" maxlength="50" show-word-limit>
            <template #prepend>tag/</template>
          </el-input>
          <div class="form-tips">将显示为: example.com/tag/{{ form.slug || 'your-slug' }}</div>
        </el-form-item>

        <el-form-item label="标签描述" prop="description">
          <el-input v-model="form.description" type="textarea" :rows="3" placeholder="请输入标签描述" maxlength="500"
            show-word-limit />
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
  name: 'TagList'
})

import { ref, reactive, watch, onMounted, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search, Refresh } from '@element-plus/icons-vue'
import {
  getTagList,
  createTag,
  updateTag,
  deleteTag
} from '@/api/tag'
import HaloButton from '@/layout/components/HaloButton.vue'

const loading = ref(false)
const dialogVisible = ref(false)
const submitting = ref(false)
// 表单引用
const formRef = ref()
const tagList = ref([])

const searchForm = reactive({
  name: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const dialogTitle = ref('新增标签')
// 表单数据
const form = reactive({
  id: null,
  name: '',
  slug: '',
  description: ''
})

// 自定义校验函数
const validateName = (rule, value, callback) => {
  if (!value || value.trim() === '') {
    callback(new Error('标签名称不能为空'))
  } else if (value.length > 50) {
    callback(new Error('标签名称长度不能超过50个字符'))
  } else {
    callback()
  }
}

const validateSlug = (rule, value, callback) => {
  if (!value || value.trim() === '') {
    callback(new Error('URL后缀不能为空'))
  } else if (value.length > 50) {
    callback(new Error('URL后缀长度不能超过50个字符'))
  } else if (!/^[a-z0-9]+(?:-[a-z0-9]+)*$/.test(value)) {
    callback(new Error('URL后缀只能包含小写字母、数字和连字符(-)，且不能以连字符开头或结尾'))
  } else {
    callback()
  }
}

const validateDescription = (rule, value, callback) => {
  if (value && value.length > 500) {
    callback(new Error('描述长度不能超过500个字符'))
  } else {
    callback()
  }
}

// 校验规则
const rules = {
  name: [
    { required: true, message: '请输入标签名称', trigger: 'blur' },
    { validator: validateName, trigger: 'blur' }
  ],
  slug: [
    { required: true, message: '请输入URL后缀', trigger: 'blur' },
    { validator: validateSlug, trigger: 'blur' }
  ],
  description: [
    { validator: validateDescription, trigger: 'blur' }
  ]
}

// 自动从名称生成slug
const autoGenerateSlug = () => {
  if (form.name && (!form.slug || form.slug === '')) {
    // 生成slug：转换为小写，移除特殊字符，空格替换为连字符
    form.slug = form.name
      .toLowerCase()
      .replace(/[^a-z0-9\s-]/g, '') // 移除非字母数字字符（除了空格和连字符）
      .replace(/\s+/g, '-')         // 空格替换为连字符
      .replace(/^-+|-+$/g, '')      // 移除开头和结尾的连字符
      .replace(/-+/g, '-')          // 多个连字符替换为单个
  }
}

// 监听名称变化自动生成slug
watch(() => form.name, () => {
  // 只有slug为空时才自动生成
  if (!form.slug || form.slug === '') {
    autoGenerateSlug()
  }
})

// 表单验证方法
const validateForm = () => {
  return new Promise((resolve, reject) => {
    if (!formRef.value) {
      reject(new Error('表单未初始化'))
      return
    }

    formRef.value.validate((valid) => {
      if (valid) {
        resolve({
          name: form.name.trim(),
          slug: form.slug.trim(),
          description: form.description ? form.description.trim() : null
        })
      } else {
        reject(new Error('请检查表单填写是否正确'))
      }
    })
  })
}

// 加载标签列表
const loadList = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      pageSize: pagination.pageSize,
      name: searchForm.name.trim() || undefined
    }

    const res = await getTagList(params)
    console.log(res.data);

    tagList.value = res.data.items || []
    pagination.total = res.data.total || 0
  } catch (error) {
    console.error('加载标签列表失败:', error)
    ElMessage.error('加载失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  loadList()
}

// 重置搜索
const handleReset = () => {
  searchForm.name = ''
  handleSearch()
}

// 新增标签
const handleCreate = () => {
  dialogTitle.value = '新增标签'
  Object.assign(form, {
    id: null,
    name: '',
    slug: '',
    description: ''
  })
  dialogVisible.value = true

  // 下次 DOM 更新后清除表单验证状态
  nextTick(() => {
    if (formRef.value) {
      formRef.value.clearValidate()
    }
  })
}

// 编辑标签
const handleEdit = (row) => {
  dialogTitle.value = '编辑标签'
  Object.assign(form, {
    id: row.id,
    name: row.name,
    slug: row.slug,
    description: row.description || ''
  })
  dialogVisible.value = true

  // 下次 DOM 更新后清除表单验证状态
  nextTick(() => {
    if (formRef.value) {
      formRef.value.clearValidate()
    }
  })
}

// 对话框关闭时的处理
const handleDialogClosed = () => {
  // 清除表单验证状态
  if (formRef.value) {
    formRef.value.clearValidate()
  }
}

// 提交表单
const handleSubmit = async () => {
  try {
    // 验证表单
    const formData = await validateForm()

    submitting.value = true

    if (form.id) {
      // 更新标签
      await updateTag(form.id, formData)
      ElMessage.success('更新成功')
    } else {
      // 创建标签
      await createTag(formData)
      ElMessage.success('创建成功')
    }

    // 关闭对话框并刷新列表
    dialogVisible.value = false
    loadList()
  } catch (error) {
    // 验证失败或API调用失败
    if (error.message !== '请检查表单填写是否正确') {
      ElMessage.error(error.message || '操作失败')
    }
  } finally {
    submitting.value = false
  }
}

// 删除标签
const handleDelete = async (id) => {
  try {
    await ElMessageBox.confirm('确定删除该标签吗？删除后无法恢复，且相关文章的标签将被移除。', '警告', {
      type: 'warning',
      confirmButtonText: '确定删除',
      cancelButtonText: '取消'
    })

    await deleteTag(id)
    ElMessage.success('删除成功')
    loadList()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 组件挂载时加载列表
onMounted(() => {
  loadList()
})
</script>

<style scoped lang="scss">
.tag-list {
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

  .form-tips {
    font-size: 12px;
    color: #909399;
    margin-top: 4px;
  }


:deep(.el-input-group__prepend) {
  background-color: #f5f7fa;
  color: #909399;
}
</style>