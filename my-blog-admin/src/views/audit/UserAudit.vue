<template>
  <div class="audit-container">
    <!-- 筛选条件 -->
    <div class="filter-section">
      <el-form :inline="true" :model="filterForm">
        <el-form-item>
          <el-input v-model="filterForm.username" placeholder="用户名" clearable style="width: 150px" />
        </el-form-item>
        <el-form-item>
          <el-input v-model="filterForm.email" placeholder="邮箱" clearable style="width: 200px" />
        </el-form-item>
        <el-form-item>
          <el-select v-model="filterForm.status" placeholder="账号状态" clearable style="width: 130px" :teleported="false">
            <el-option label="待审核" value="0" />
            <el-option label="正常" value="1" />
            <el-option label="禁用" value="2" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">查询</el-button>
          <el-button @click="resetFilter">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 用户列表 -->
    <div class="table-section">
      <el-table :data="userList" v-loading="loading">
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="username" label="用户名" width="110" />
        <el-table-column prop="nickname" label="昵称" width="110" />
        <el-table-column prop="email" label="邮箱" min-width="180" />
        <el-table-column prop="role" label="角色" width="80" />
        <el-table-column prop="createdAt" label="注册时间" width="160" />
        <el-table-column label="状态" width="90" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.status === 0" type="warning" size="small">待审核</el-tag>
            <el-tag v-else-if="row.status === 1" type="success" size="small">正常</el-tag>
            <el-tag v-else-if="row.status === 2" type="danger" size="small">禁用</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right" align="center">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="viewUser(row)">查看</el-button>
            <el-button v-if="row.status === 0" link type="success" size="small" @click="approveUser(row)">通过</el-button>
            <el-button v-if="row.status === 0" link type="danger" size="small" @click="rejectUser(row)">拒绝</el-button>
            <el-button v-if="row.status === 1" link type="warning" size="small" @click="disableUser(row)">禁用</el-button>
            <el-button v-if="row.status === 2" link type="success" size="small" @click="enableUser(row)">启用</el-button>
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

    <!-- 用户详情对话框 -->
    <el-dialog
      v-model="detailDialogVisible"
      title="用户详情"
      width="700px"
    >
      <div v-if="currentUser" class="user-detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="ID">{{ currentUser.id }}</el-descriptions-item>
          <el-descriptions-item label="用户名">{{ currentUser.username }}</el-descriptions-item>
          <el-descriptions-item label="昵称">{{ currentUser.nickname || '-' }}</el-descriptions-item>
          <el-descriptions-item label="邮箱">{{ currentUser.email }}</el-descriptions-item>
          <el-descriptions-item label="角色">{{ currentUser.role }}</el-descriptions-item>
          <el-descriptions-item label="账号状态">
            <el-tag v-if="currentUser.status === 0" type="warning">待审核</el-tag>
            <el-tag v-else-if="currentUser.status === 1" type="success">正常</el-tag>
            <el-tag v-else-if="currentUser.status === 2" type="danger">禁用</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="注册时间" :span="2">{{ currentUser.createdAt }}</el-descriptions-item>
        </el-descriptions>
        <div class="user-extra" v-if="currentUser.bio || currentUser.website || currentUser.github">
          <h4>其他信息</h4>
          <div v-if="currentUser.bio">
            <span class="label">个人简介：</span>
            <span class="value">{{ currentUser.bio }}</span>
          </div>
          <div v-if="currentUser.website">
            <span class="label">个人网站：</span>
            <a :href="currentUser.website" target="_blank">{{ currentUser.website }}</a>
          </div>
          <div v-if="currentUser.github">
            <span class="label">GitHub：</span>
            <a :href="currentUser.github" target="_blank">{{ currentUser.github }}</a>
          </div>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
defineOptions({
  name: 'UserAudit'
})

import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'

const loading = ref(false)
const submitting = ref(false)
const userList = ref([])
const auditDialogVisible = ref(false)
const detailDialogVisible = ref(false)
const currentUser = ref(null)
const auditType = ref('')
const dialogTitle = ref('')

const filterForm = reactive({
  username: '',
  email: '',
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

// 获取用户列表
const fetchUsers = async () => {
  loading.value = true
  try {
    // TODO: 调用后端接口
    // const res = await api.getPendingUsers({
    //   page: pagination.page,
    //   pageSize: pagination.pageSize,
    //   ...filterForm
    // })
    // userList.value = res.data.items
    // pagination.total = res.data.total

    // Mock 数据
    userList.value = [
      {
        id: 1,
        username: 'user001',
        nickname: '张三',
        email: 'zhangsan@example.com',
        role: '普通用户',
        status: 0,
        createdAt: '2024-01-15 10:30:00',
        bio: '热爱技术，乐于分享',
        website: 'https://zhangsan.com',
        github: 'https://github.com/zhangsan'
      },
      {
        id: 2,
        username: 'user002',
        nickname: '李四',
        email: 'lisi@example.com',
        role: '普通用户',
        status: 1,
        createdAt: '2024-01-14 14:20:00'
      },
      {
        id: 3,
        username: 'user003',
        nickname: '王五',
        email: 'wangwu@example.com',
        role: '普通用户',
        status: 0,
        createdAt: '2024-01-13 09:15:00'
      }
    ]
    pagination.total = 3
  } catch (error) {
    ElMessage.error('获取用户列表失败')
    console.error(error)
  } finally {
    loading.value = false
  }
}

// 刷新列表
const refreshList = () => {
  fetchUsers()
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchUsers()
}

// 重置筛选
const resetFilter = () => {
  filterForm.username = ''
  filterForm.email = ''
  filterForm.status = ''
  pagination.page = 1
  fetchUsers()
}

// 查看用户
const viewUser = (user) => {
  currentUser.value = user
  detailDialogVisible.value = true
}

// 审核用户
const approveUser = (user) => {
  currentUser.value = user
  auditType.value = 'approve'
  dialogTitle.value = '审核通过'
  auditForm.reason = ''
  auditDialogVisible.value = true
}

// 拒绝用户
const rejectUser = (user) => {
  currentUser.value = user
  auditType.value = 'reject'
  dialogTitle.value = '审核拒绝'
  auditForm.reason = ''
  auditDialogVisible.value = true
}

// 禁用用户
const disableUser = async (user) => {
  try {
    await ElMessageBox.confirm(
      `确定要禁用用户 "${user.username}" 吗？`,
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    // TODO: 调用后端接口
    // await api.updateUserStatus(user.id, 2)
    ElMessage.success('禁用成功')
    refreshList()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('禁用失败')
      console.error(error)
    }
  }
}

// 启用用户
const enableUser = async (user) => {
  try {
    await ElMessageBox.confirm(
      `确定要启用用户 "${user.username}" 吗？`,
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    // TODO: 调用后端接口
    // await api.updateUserStatus(user.id, 1)
    ElMessage.success('启用成功')
    refreshList()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('启用失败')
      console.error(error)
    }
  }
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
    //   await api.approveUser(currentUser.value.id, auditForm.reason)
    // } else {
    //   await api.rejectUser(currentUser.value.id, auditForm.reason)
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
  currentUser.value = null
}

onMounted(() => {
  fetchUsers()
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

  .user-detail {
    .user-extra {
      margin-top: 20px;

      h4 {
        margin-bottom: 15px;
        font-size: 14px;
        font-weight: 600;
        color: #333;
      }

      > div {
        margin-bottom: 10px;
        font-size: 13px;
        color: #666;

        .label {
          display: inline-block;
          width: 80px;
          font-weight: 500;
          color: #333;
        }

        .value {
          color: #666;
        }

        a {
          color: #1890ff;
          text-decoration: none;

          &:hover {
            text-decoration: underline;
          }
        }
      }
    }
  }
}
</style>
