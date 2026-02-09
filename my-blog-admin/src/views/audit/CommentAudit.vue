<template>
  <div class="audit-container">
    <!-- 筛选条件 -->
    <div class="filter-section">
      <el-form :inline="true" :model="filterForm">
        <el-form-item>
          <el-input v-model="filterForm.content" placeholder="评论内容" clearable style="width: 200px" />
        </el-form-item>
        <el-form-item>
          <el-input v-model="filterForm.user" placeholder="评论者" clearable style="width: 120px" />
        </el-form-item>
        <el-form-item>
          <el-select v-model="filterForm.status" placeholder="审核状态" clearable style="width: 130px" :teleported="false">
            <el-option label="待审核" value="0" />
            <el-option label="已通过" value="1" />
            <el-option label="已拒绝" value="2" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-select v-model="filterForm.level" placeholder="评论层级" clearable style="width: 130px">
            <el-option label="一级评论" value="0" />
            <el-option label="二级评论" value="1" />
            <el-option label="三级评论" value="2" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">查询</el-button>
          <el-button @click="resetFilter">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 评论列表 -->
    <div class="table-section">
      <el-table :data="commentList" v-loading="loading" :row-class-name="getRowClassName">
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column label="层级" width="80" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.level === 0" size="small" type="info">一级</el-tag>
            <el-tag v-else-if="row.level === 1" size="small" type="primary">二级</el-tag>
            <el-tag v-else-if="row.level === 2" size="small" type="success">三级</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="评论内容" min-width="260">
          <template #default="{ row }">
            <div class="comment-content-cell">
              <div class="content-text">{{ row.content }}</div>
              <div v-if="row.parentUser" class="parent-info">回复 {{ row.parentUser }}</div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="article" label="所属文章" width="140" show-overflow-tooltip />
        <el-table-column prop="user" label="评论者" width="100" />
        <el-table-column prop="createdAt" label="评论时间" width="160" />
        <el-table-column label="状态" width="90" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.status === 0" type="warning" size="small">待审核</el-tag>
            <el-tag v-else-if="row.status === 1" type="success" size="small">已通过</el-tag>
            <el-tag v-else-if="row.status === 2" type="danger" size="small">已拒绝</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right" align="center">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="viewComment(row)">查看</el-button>
            <el-button v-if="row.status === 0" link type="success" size="small" @click="approveComment(row)">通过</el-button>
            <el-button v-if="row.status === 0" link type="danger" size="small" @click="rejectComment(row)">拒绝</el-button>
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

    <!-- 评论详情对话框 -->
    <el-dialog
      v-model="detailDialogVisible"
      title="评论详情"
      width="700px"
    >
      <div v-if="currentComment" class="comment-detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="ID">{{ currentComment.id }}</el-descriptions-item>
          <el-descriptions-item label="评论者">{{ currentComment.user }}</el-descriptions-item>
          <el-descriptions-item label="所属文章" :span="2">{{ currentComment.article }}</el-descriptions-item>
          <el-descriptions-item label="评论时间">{{ currentComment.createdAt }}</el-descriptions-item>
          <el-descriptions-item label="审核状态">
            <el-tag v-if="currentComment.status === 0" type="warning">待审核</el-tag>
            <el-tag v-else-if="currentComment.status === 1" type="success">已通过</el-tag>
            <el-tag v-else-if="currentComment.status === 2" type="danger">已拒绝</el-tag>
          </el-descriptions-item>
        </el-descriptions>
        <div class="comment-content">
          <h4>评论内容</h4>
          <p>{{ currentComment.content }}</p>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
defineOptions({
  name: 'CommentAudit'
})

import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'

const loading = ref(false)
const submitting = ref(false)
const commentList = ref([])
const auditDialogVisible = ref(false)
const detailDialogVisible = ref(false)
const currentComment = ref(null)
const auditType = ref('')
const dialogTitle = ref('')

const filterForm = reactive({
  content: '',
  user: '',
  status: '',
  level: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const auditForm = reactive({
  reason: ''
})

// 获取评论列表
const fetchComments = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      pageSize: pagination.pageSize
    }
    if (filterForm.content) params.content = filterForm.content
    if (filterForm.user) params.user = filterForm.user
    if (filterForm.status !== '') params.status = filterForm.status
    if (filterForm.level !== '') params.level = filterForm.level

    const response = await fetch('/api/v1/admin/comments/pending?' + new URLSearchParams(params), {
      headers: {
        'Authorization': 'Bearer ' + localStorage.getItem('token')
      }
    })

    const data = await response.json()

    if (data.code === 200) {
      // 转换数据格式，添加层级信息和父用户名
      commentList.value = data.data.items.map(item => {
        const level = item.parent_id ? 1 : 0
        return {
          id: item.id,
          content: item.content,
          article: item.article?.title || '未知文章',
          articleId: item.article_id,
          user: item.user?.username || item.user?.email || '未知用户',
          createdAt: item.created_at,
          status: item.status,
          level: level,
          parentUser: item.parent?.user?.username || null
        }
      })
      pagination.total = data.data.total
    } else {
      ElMessage.error(data.message || '获取评论列表失败')
    }
  } catch (error) {
    ElMessage.error('获取评论列表失败')
    console.error(error)
  } finally {
    loading.value = false
  }
}

// 刷新列表
const refreshList = () => {
  fetchComments()
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchComments()
}

// 重置筛选
const resetFilter = () => {
  filterForm.content = ''
  filterForm.user = ''
  filterForm.status = ''
  filterForm.level = ''
  pagination.page = 1
  fetchComments()
}

// 获取行样式
const getRowClassName = ({ row }) => {
  if (row.level === 0) return 'level-0'
  if (row.level === 1) return 'level-1'
  if (row.level === 2) return 'level-2'
  return ''
}

// 查看评论
const viewComment = (comment) => {
  currentComment.value = comment
  detailDialogVisible.value = true
}

// 审核评论
const approveComment = (comment) => {
  currentComment.value = comment
  auditType.value = 'approve'
  dialogTitle.value = '审核通过'
  auditForm.reason = ''
  auditDialogVisible.value = true
}

// 拒绝评论
const rejectComment = (comment) => {
  currentComment.value = comment
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
    const endpoint = auditType.value === 'approve'
      ? `/api/v1/admin/comments/${currentComment.value.id}/approve`
      : `/api/v1/admin/comments/${currentComment.value.id}/reject`

    const response = await fetch(endpoint, {
      method: 'PUT',
      headers: {
        'Authorization': 'Bearer ' + localStorage.getItem('token'),
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ reason: auditForm.reason })
    })

    const data = await response.json()

    if (data.code === 200) {
      ElMessage.success('审核成功')
      auditDialogVisible.value = false
      refreshList()
    } else {
      ElMessage.error(data.message || '审核失败')
    }
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
  currentComment.value = null
}

onMounted(() => {
  fetchComments()
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

      &.level-1 td {
        background-color: #f8fbfd;
      }

      &.level-2 td {
        background-color: #f0f7fc;
      }
    }

    .comment-content-cell {
      .content-text {
        color: #333;
        line-height: 1.5;
        margin-bottom: 4px;
      }

      .parent-info {
        font-size: 12px;
        color: #909399;
        display: flex;
        align-items: center;

        &::before {
          content: '';
          display: inline-block;
          width: 0;
          height: 0;
          border-left: 4px solid transparent;
          border-right: 4px solid transparent;
          border-top: 4px solid #909399;
          margin-right: 6px;
        }
      }
    }

    .pagination-container {
      display: flex;
      justify-content: flex-end;
      margin-top: 20px;
      padding-top: 16px;
      border-top: 1px solid #f0f0f0;
    }
  }

  .comment-detail {
    .comment-content {
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
