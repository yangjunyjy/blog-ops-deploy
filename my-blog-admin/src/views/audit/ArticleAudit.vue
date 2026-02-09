<template>
  <div class="audit-container">
    <!-- 筛选条件 -->
    <div class="filter-section">
      <el-form :inline="true" :model="filterForm">
        <el-form-item>
          <el-input v-model="filterForm.title" placeholder="文章标题" clearable style="width: 200px" prefix-icon="Search" />
        </el-form-item>
        <el-form-item>
          <el-input v-model="filterForm.author" placeholder="作者" clearable style="width: 150px" />
        </el-form-item>
        <el-form-item>
          <el-select v-model="filterForm.status" placeholder="审核状态" clearable style="width: 130px" :teleported="false">
            <el-option label="待审核" value="0" />
            <el-option label="已通过" value="1" />
            <el-option label="已拒绝" value="2" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">查询</el-button>
          <el-button @click="resetFilter">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 文章列表 -->
    <div class="table-section">
      <el-table :data="articleList" v-loading="loading">
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="title" label="标题" min-width="240" show-overflow-tooltip />
        <el-table-column prop="author" label="作者" width="100" />
        <el-table-column prop="category" label="分类" width="90" />
        <el-table-column prop="createdAt" label="创建时间" width="160" />
        <el-table-column label="状态" width="90" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.status === 0" type="warning" size="small">待审核</el-tag>
            <el-tag v-else-if="row.status === 1" type="success" size="small">已通过</el-tag>
            <el-tag v-else-if="row.status === 2" type="danger" size="small">已拒绝</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right" align="center">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="viewArticle(row)">查看</el-button>
            <el-button v-if="row.status === 0" link type="success" size="small" @click="approveArticle(row)">通过</el-button>
            <el-button v-if="row.status === 0" link type="danger" size="small" @click="rejectArticle(row)">拒绝</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :total="pagination.total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next"
          @size-change="refreshList"
          @current-change="refreshList"
        />
      </div>
    </div>

    <!-- 审核对话框 -->
    <el-dialog
      v-model="auditDialogVisible"
      :title="dialogTitle"
      width="500px"
      @close="resetAuditForm"
    >
      <el-form :model="auditForm" label-width="80px">
        <el-form-item label="拒绝原因" v-if="auditType === 'reject'">
          <el-input
            v-model="auditForm.reason"
            type="textarea"
            :rows="4"
            placeholder="请输入拒绝原因"
          />
        </el-form-item>
        <el-form-item label="审核说明" v-if="auditType === 'approve'">
          <el-input
            v-model="auditForm.reason"
            type="textarea"
            :rows="4"
            placeholder="请输入审核说明（可选）"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="auditDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmAudit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>

    <!-- 文章详情对话框 -->
    <el-dialog
      v-model="detailDialogVisible"
      title="文章详情"
      width="800px"
    >
      <div v-if="currentArticle" class="article-detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="ID">{{ currentArticle.id }}</el-descriptions-item>
          <el-descriptions-item label="标题">{{ currentArticle.title }}</el-descriptions-item>
          <el-descriptions-item label="作者">{{ currentArticle.author }}</el-descriptions-item>
          <el-descriptions-item label="分类">{{ currentArticle.category }}</el-descriptions-item>
          <el-descriptions-item label="创建时间">{{ currentArticle.createdAt }}</el-descriptions-item>
          <el-descriptions-item label="审核状态">
            <el-tag v-if="currentArticle.status === 0" type="warning">待审核</el-tag>
            <el-tag v-else-if="currentArticle.status === 1" type="success">已通过</el-tag>
            <el-tag v-else-if="currentArticle.status === 2" type="danger">已拒绝</el-tag>
          </el-descriptions-item>
        </el-descriptions>
        <div class="article-content">
          <h4>文章摘要</h4>
          <p>{{ currentArticle.summary }}</p>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
defineOptions({
  name: 'ArticleAudit'
})

import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'

const loading = ref(false)
const submitting = ref(false)
const articleList = ref([])
const auditDialogVisible = ref(false)
const detailDialogVisible = ref(false)
const currentArticle = ref(null)
const auditType = ref('')
const dialogTitle = ref('')

const filterForm = reactive({
  title: '',
  author: '',
  status: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const auditForm = reactive({
  reason: ''
})

// 获取文章列表
const fetchArticles = async () => {
  loading.value = true
  try {
    // TODO: 调用后端接口
    // const res = await api.getPendingArticles({
    //   page: pagination.page,
    //   pageSize: pagination.pageSize,
    //   ...filterForm
    // })
    // articleList.value = res.data.items
    // pagination.total = res.data.total

    // Mock 数据
    articleList.value = [
      {
        id: 1,
        title: 'Vue 3 组合式 API 最佳实践',
        author: '张三',
        category: '前端开发',
        createdAt: '2024-01-15 10:30:00',
        status: 0,
        summary: '本文介绍 Vue 3 组合式 API 的使用方法和最佳实践...'
      },
      {
        id: 2,
        title: 'Go 语言并发编程指南',
        author: '李四',
        category: '后端开发',
        createdAt: '2024-01-14 14:20:00',
        status: 0,
        summary: '深入理解 Go 语言的 goroutine 和 channel...'
      },
      {
        id: 3,
        title: 'Docker 容器化部署',
        author: '王五',
        category: 'DevOps',
        createdAt: '2024-01-13 09:15:00',
        status: 1,
        summary: '使用 Docker 容器化部署应用程序...'
      }
    ]
    pagination.total = 3
  } catch (error) {
    ElMessage.error('获取文章列表失败')
    console.error(error)
  } finally {
    loading.value = false
  }
}

// 刷新列表
const refreshList = () => {
  fetchArticles()
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchArticles()
}

// 重置筛选
const resetFilter = () => {
  filterForm.title = ''
  filterForm.author = ''
  filterForm.status = ''
  pagination.page = 1
  fetchArticles()
}

// 查看文章
const viewArticle = (article) => {
  currentArticle.value = article
  detailDialogVisible.value = true
}

// 审核文章
const approveArticle = (article) => {
  currentArticle.value = article
  auditType.value = 'approve'
  dialogTitle.value = '审核通过'
  auditForm.reason = ''
  auditDialogVisible.value = true
}

// 拒绝文章
const rejectArticle = (article) => {
  currentArticle.value = article
  auditType.value = 'reject'
  dialogTitle.value = '审核拒绝'
  auditForm.reason = ''
  auditDialogVisible.value = true
}

// 确认审核
const confirmAudit = async () => {
  if (auditType.value === 'reject' && !auditForm.reason.trim()) {
    ElMessage.warning('请输入拒绝原因')
    return
  }

  submitting.value = true
  try {
    // TODO: 调用后端接口
    // if (auditType.value === 'approve') {
    //   await api.approveArticle(currentArticle.value.id, auditForm.reason)
    // } else {
    //   await api.rejectArticle(currentArticle.value.id, auditForm.reason)
    // }

    ElMessage.success('审核成功')
    auditDialogVisible.value = false
    refreshList()
  } catch (error) {
    ElMessage.error('审核失败')
    console.error(error)
  } finally {
    submitting.value = false
  }
}

// 重置审核表单
const resetAuditForm = () => {
  auditForm.reason = ''
  currentArticle.value = null
}

onMounted(() => {
  fetchArticles()
})
</script>

<style scoped lang="scss">
.audit-container {
  height: 100%;
  display: flex;
  flex-direction: column;

  .filter-section {
    background: #fff;
    padding: 16px 20px;
    border-radius: 4px;
    margin-bottom: 12px;
    flex-shrink: 0;
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);

    :deep(.el-form-item) {
      margin-bottom: 0;
    }
  }

  .table-section {
    background: #fff;
    padding: 16px;
    border-radius: 4px;
    flex: 1;
    display: flex;
    flex-direction: column;
    min-height: 0;
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);

    :deep(.el-table) {
      font-size: 13px;
      flex: 1;

      .el-table__header th {
        background: #fafafa;
        font-weight: 500;
        color: #333;
      }
    }

    .pagination-container {
      display: flex;
      justify-content: flex-end;
      margin-top: 12px;
      padding-top: 12px;
      border-top: 1px solid #f0f0f0;
      flex-shrink: 0;
    }
  }

  .article-detail {
    .article-content {
      margin-top: 20px;

      h4 {
        margin-bottom: 10px;
        font-size: 14px;
        font-weight: 600;
        color: #333;
      }

      p {
        margin: 0;
        color: #666;
        line-height: 1.6;
      }
    }
  }
}
</style>
