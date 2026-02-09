<template>
  <div class="user-list">
    <el-card shadow="never">
      <!-- 搜索栏 -->
      <el-form :model="searchForm" inline class="search-form" @submit.prevent>
        <div class="form-items">
          <el-form-item label="用户名">
            <el-input v-model="searchForm.username" placeholder="用户名" clearable style="width: 200px" />
          </el-form-item>
          <el-form-item label="邮箱">
            <el-input v-model="searchForm.email" placeholder="邮箱" clearable style="width: 200px" />
          </el-form-item>
        </div>
        <el-form-item>
          <HaloButton type="primary" size="small" content="查询" @click.prevent="handleSearch" />
          <HaloButton type="default" size="small" content="重置" @click.prevent="handleReset" />
        </el-form-item>
      </el-form>

      <!-- 表格 -->
      <div class="table-wrapper">
        <!-- 操作栏 -->
        <div class="action-bar">
          <span class="action-title">用户管理</span>
          <div class="action-buttons">
            <HaloButton type="primary" size="medium" :icon="Plus" content="新增" @click="handleCreate" />
          </div>
        </div>

        <el-table :data="userList" v-loading="loading" :header-cell-style="{ backgroundColor: '#f5f7fa' }">
          <el-table-column prop="id" label="ID" width="70" align="center" />
          <el-table-column prop="avatar" label="头像" width="100">
            <template #default="{ row }">
              <el-avatar :size="40" :src="row.avatar" />
            </template>
          </el-table-column>
          <el-table-column prop="username" label="用户名" min-width="120" />
          <el-table-column prop="nickname" label="昵称" min-width="120" />
          <el-table-column prop="email" label="邮箱" min-width="180" />
          <el-table-column prop="roles" label="角色" min-width="150">
            <template #default="{ row }">
              <el-tag v-for="role in row.roles" :key="role.id" :type="getRoleTagType(role.roleCode)" size="small"
                style="margin-right: 5px">
                {{ role.roleName }}
              </el-tag>
              <span v-if="!row.roles || row.roles.length === 0" style="color: #999">-</span>
            </template>
          </el-table-column>
          <el-table-column label="状态" width="80">
            <template #default="{ row }">
              <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
                {{ row.status === 1 ? '启用' : '禁用' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="创建时间" width="160">
            <template #default="{ row }">
              {{ formatDate(row.createTime) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="140" fixed="right">
            <template #default="{ row }">
              <div style="display: flex; align-items: center; gap: 4px;">
                <el-button type="primary" link size="small" @click="handleEdit(row)" class="icon-btn" v-if="isAdmin()">
                  <el-icon class="action-icon">
                    <img src="/svg/user_edit.svg" alt="edit" style="width: 18px; height: 18px;" />
                  </el-icon>
                </el-button>
                <el-button type="success" link size="small" @click="handleAssignRoles(row)" class="icon-btn">
                  <el-icon class="action-icon">
                    <img src="/svg/user_assign.svg" alt="assign" style="width: 18px; height: 18px;" />
                  </el-icon>
                </el-button>
                <el-button type="danger" link size="small" @click="handleDelete(row.id)" class="icon-btn"
                  v-if="isAdmin()">
                  <el-icon class="action-icon delete-icon">
                    <img src="/svg/user_delete.svg" alt="delete" style="width: 18px; height: 18px;" />
                  </el-icon>
                </el-button>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <!-- 分页器 -->
      <div class="pagination-wrapper">
        <el-pagination v-model:current-page="pagination.page" v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50, 100]" :total="pagination.total" layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange" @current-change="handleCurrentChange" />
      </div>
    </el-card>

    <!-- 编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="form.username" placeholder="请输入用户名" :disabled="!!form.id" />
        </el-form-item>
        <el-form-item label="昵称" prop="nickname">
          <el-input v-model="form.nickname" placeholder="请输入用户名" :disabled="!!form.id" />
        </el-form-item>
        <el-form-item label="密码" prop="password" type="password">
          <el-input v-model="form.password" placeholder="请输入密码" :disabled="!!form.id" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="form.email" placeholder="请输入邮箱" />
        </el-form-item>
        <el-form-item label="性别" prop="gender">
          <el-radio-group v-model="form.gender">
            <el-radio :value="0">男</el-radio>
            <el-radio :value="1">女</el-radio>
            <el-radio :value="2">未知</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio :value="1">启用</el-radio>
            <el-radio :value="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>

    <!-- 分配角色对话框 -->
    <RoleAssignDialog v-model="assignDialogVisible" :user-id="currentUser?.id"
      :existing-role-ids="currentUser?.roles || []" @confirm="handleAssignConfirm" />
  </div>
</template>

<script setup>
defineOptions({
  name: 'UserList'
})

import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import HaloButton from '../../layout/components/HaloButton.vue'
import RoleAssignDialog from './components/roleAssign'
import { useRbacStore } from '../../store/rbac'
import { getUserList, assignUserRoles, deleteUser, updateUser, createUser } from '@/api/user'

const loading = ref(false)
const dialogVisible = ref(false)
const assignDialogVisible = ref(false)
const submitting = ref(false)
const formRef = ref(null)
const userList = ref([])
const currentUser = ref(null)
const currentLoginUser = ref(null)
const rbacStore = useRbacStore()


// 分页
currentLoginUser.value = rbacStore.currentUser
const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})




const searchForm = reactive({
  username: '',
  email: '',
  status: 1
})

const dialogTitle = ref('新增用户')
const form = reactive({
  id: null,
  nickname: '',
  username: '',
  password: '',
  gender: 0,
  email: '',
  status: 1
})

const isAdmin = () => {
  return currentLoginUser.value.isAdmin === 1
}
const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  nickname: [{ required: true, message: '请输入昵称', trigger: 'blur' }],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ]
}

const formatDate = (date) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}

const filterRoleID = (user) => {
  const roleIds = user.roles.map(role => role.id)
  user.roleIds = roleIds
}

const getRoleTagType = (roleCode) => {
  if (roleCode == "admin") {
    return 'danger'
  }
  return 'info'
}

const loadList = async () => {
  const params = {
    username: searchForm.username,
    email: searchForm.email,
    status: searchForm.status,
    page: pagination.page,
    pageSize: pagination.pageSize
  }
  const data = await getUserList(params)

  console.log("请求的用户数据为", data.data);

  if (data.data.list && data.data.list.length > 0) {
    console.log("用户列表第一条:", JSON.stringify(data.data.list[0], null, 2));
    console.log("第一条用户的角色:", data.data.list[0].roles);
    if (data.data.list[0].roles && data.data.list[0].roles.length > 0) {
      console.log("第一个角色对象:", JSON.stringify(data.data.list[0].roles[0], null, 2));
    }
  }

  console.log("当前用户信息", currentLoginUser.value);

  userList.value = data.data.list
  userList.value.forEach(user => {
    filterRoleID(user)
  })
  pagination.total = data.data.total
}

const handleSearch = () => {
  pagination.page = 1
  loadList()
}

const handleReset = () => {
  searchForm.username = ''
  searchForm.email = ''
  pagination.page = 1
  loadList()
}

const handleCreate = () => {
  dialogTitle.value = '新增用户'
  Object.assign(form, {
    id: null,
    username: '',
    nickname: '',
    password: '',
    gender: 0,
    email: '',
    status: 1
  })
  dialogVisible.value = true
}

const handleEdit = (row) => {
  dialogTitle.value = '编辑用户'
  Object.assign(form, row)
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitting.value = true
    try {
      if (form.id) {
        const updateParams = {
          id: form.id,
          nickname: form.nickname,
          email: form.email,
          status: Number(form.status),
          gender: Number(form.gender),
          username: form.username,
          password: form.password
        }
        const data = await updateUser(updateParams)
        if (data.code != 200) {
          ElMessage.error("更新用户失败")
          return
        }
        ElMessage.success('更新成功')
      } else {
        const createParams = {
          username: form.username,
          nickname: form.nickname,
          password: form.password,
          email: form.email,
          status: Number(form.status),
          gender: Number(form.gender)
        }
        console.log('创建用户参数:', JSON.stringify(createParams, null, 2))
        const data = await createUser(createParams)
        if (data.code != 200) {
          ElMessage.error("创建用户失败")
          submitting.value = false
          return
        }
        ElMessage.success('创建用户成功')
      }
      dialogVisible.value = false
      loadList()
    } catch (error) {
    } finally {
      submitting.value = false
    }
  })
}

const handleAssignRoles = (row) => {
  console.log('handleAssignRoles - row 完整对象:', JSON.stringify(row, null, 2))
  console.log('handleAssignRoles - row.id:', row.id)
  console.log('handleAssignRoles - row.ID:', row.ID)
  console.log('handleAssignRoles - row.Id:', row.Id)
  currentUser.value = row
  assignDialogVisible.value = true
}

const handleAssignConfirm = async (data) => {
  try {
    console.log('handleAssignConfirm 接收到的数据:', data)
    const assignParams = {
      userId: data.userId || data.user_id || data.userID,
      roleIds: data.roleIds || data.role_ids || data.role_IDs
    }
    console.log('发送到后端的参数:', assignParams)
    const value = await assignUserRoles(assignParams)
    if (value.code != 200) {
      ElMessage.error('分配角色失败')
    }
    ElMessage.success("分配角色成功")
    // 如果分配的是当前登录用户，刷新权限和菜单
    const userInfo = JSON.parse(sessionStorage.getItem('userInfo') || '{}')
    if (currentUser.value.email === userInfo.email) {
      rbacStore.refreshUserPermission()
    }
    loadList()
  } catch (error) {
    ElMessage.error('分配角色失败')
  }
}

const handleDelete = async (id) => {
  await ElMessageBox.confirm('确定删除该用户吗？', '提示', { type: 'warning' })
  const data = await deleteUser(id)
  if (data.code != 200) {
    ElMessage.success('删除失败')
    return
  }
  ElMessage.success('删除成功')
  loadList()
}

const handleSizeChange = (size) => {
  pagination.pageSize = size
  pagination.page = 1
}

const handleCurrentChange = (page) => {
  pagination.page = page
}

onMounted(() => {
  loadList()
})
</script>

<style scoped lang="scss">
.user-list {
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
    flex-shrink: 0;
  }

  .table-wrapper {
    border: 1px solid #ebeef5;
    border-radius: 4px;
    padding: 16px;
    flex: 1;
    display: flex;
    flex-direction: column;
    min-height: 0;

    .action-bar {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 12px;
      flex-shrink: 0;

      .action-title {
        font-size: 16px;
        font-weight: 600;
        color: #333539;
      }

      .action-buttons {
        display: flex;
        gap: 2px;
      }
    }

    :deep(.el-table) {
      flex: 1;
      overflow: auto;
    }

    :deep(.el-table__body-wrapper) {
      overflow: auto;
    }
  }

  .pagination-wrapper {
    display: flex;
    justify-content: flex-end;
    padding-top: 12px;
    flex-shrink: 0;
  }

  :deep(.el-table__header) {
    background-color: #f5f7fa !important;
  }
}

.search-form {
  border: 1px solid #ebeef5;
  border-radius: 4px;
  padding: 16px 20px 16px 20px;
  display: flex;
  align-items: center;
  justify-content: space-between;

  .form-items {
    display: flex;
    align-items: center;
    gap: 20px;
  }

  :deep(.el-form-item:last-child) {
    margin-right: 0;
  }

  :deep(.el-form-item) {
    margin-bottom: 0;
    margin-top: 0;
  }
}

.action-icon {
  font-size: 16px;
  padding: 4px;
  border-radius: 4px;
  transition: all 0.3s;

  &:hover {
    transform: scale(1.1);
    background-color: rgba(64, 158, 255, 0.1);
  }

  &.delete-icon {
    img {
      filter: brightness(0) saturate(100%) invert(22%) sepia(87%) saturate(3827%) hue-rotate(352deg) brightness(92%) contrast(85%);
    }
  }
}

.icon-btn {
  margin: 0 2px;
  padding: 4px;

  &:hover .action-icon {
    background-color: rgba(64, 158, 255, 0.1);
  }
}
</style>
