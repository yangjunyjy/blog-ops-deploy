<template>
  <el-dialog v-model="visible" :title="dialogTitle" width="900px" @open="loadRoleList()" @close="resetRoleList()"
    :close-on-click-modal="false" :destroy-on-close="true" class="role-assign-dialog">
    <!-- 中间区域 -->
    <div class="dialog-content">
      <!-- 左侧区域 -->
      <div class="left-section">
        <!-- 左上：搜索区域 -->
        <div class="search-area">
          <div class="search-form">
            <el-form :inline="true" :model="searchForm" class="search-inline-form" @submit.prevent>
              <el-form-item label="角色名称">
                <el-input v-model="searchForm.roleName" placeholder="角色名称" clearable class="form-input" />
              </el-form-item>
              <el-form-item label="角色编码">
                <el-input v-model="searchForm.roleCode" placeholder="角色编码" clearable class="form-input" />
              </el-form-item>
              <el-form-item>
                <el-button type="primary" :icon="Search" @click.prevent="handleSearch" size="default">
                  查询
                </el-button>
              </el-form-item>
            </el-form>
          </div>
        </div>

        <!-- 左中：角色列表 -->
        <div class="role-list-area">
          <el-table ref="roleTableRef" :data="filteredRoles" v-loading="loading" height="320"
            @selection-change="handleSelectionChange" class="role-table">
            <el-table-column type="selection" width="55" align="center" />
            <el-table-column prop="roleName" label="角色名称" min-width="150" />
            <el-table-column prop="roleDesc" label="角色描述" min-width="200" show-overflow-tooltip />
          </el-table>
        </div>

        <!-- 左下：分页器 -->
        <div class="pagination-area">
          <el-pagination v-model:current-page="pagination.page" v-model:page-size="pagination.pageSize"
            :total="pagination.total" :page-sizes="[10, 20, 50]" size="small" layout="total, prev, pager, next, jumper"
            @size-change="handleSizeChange" @current-change="handleCurrentChange" />
        </div>
      </div>

      <!-- 右侧区域：已选角色 -->
      <div class="right-section">
        <div class="selected-header">
          <span class="selected-title">已选角色</span>
          <span class="selected-count">({{ selectedRoles.length }})</span>
        </div>
        <div class="selected-list">
          <el-empty v-if="selectedRoles.length === 0" description="暂未选择角色" :image-size="80" />
          <div v-else class="selected-items">
            <div v-for="role in selectedRoles" :key="role.id" class="selected-item">
              <div class="role-info">
                <el-tag :type="getRoleTagType(role)" size="default">
                  {{ role.roleName }}
                </el-tag>
                <span class="role-code">{{ role.roleCode }}</span>
              </div>
              <el-button type="danger" link size="small" :icon="Close" @click="handleRemoveRole(role.id)">
                移除
              </el-button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 底部按钮 -->
    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleCancel" size="default">取消</el-button>
        <el-button type="primary" @click="handleConfirm" :loading="submitting" size="default">
          确认
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { Search, Close } from '@element-plus/icons-vue'
import { getRoleList } from '@/api/role'

// Props
const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  title: {
    type: String,
    default: '分配角色'
  },
  userId: {
    type: [Number, String],
    default: null
  },
  userName: {
    type: String,
    default: ''
  },
  existingRoleIds: {
    type: Array,
    default: () => []
  }
})

// Emits
const emit = defineEmits(['update:modelValue', 'confirm', 'cancel'])

// 响应式数据
const visible = ref(false)
const searchForm = ref({
  roleName: '',
  roleCode: ''
})
const loading = ref(false)
const submitting = ref(false)
const roleTableRef = ref(null)
const tableSelection = ref([])
const selectedRoles = ref([])

const roleIds = ref([])


// 分页数据
const pagination = ref({
  page: 1,
  pageSize: 10,
  total: 0
})

// 所有角色数据（实际项目中应该从API获取）
const allRoles = ref([])

const translate = () => {
  console.log('translate - existingRoleIds:', props.existingRoleIds)
  roleIds.value = props.existingRoleIds.map(role => role.id || role)
  console.log('translate - roleIds:', roleIds.value)
}
const loadRoleList = async () => {
  const dataParams = {
    roleName: searchForm.value.roleName,
    roleCode: searchForm.value.roleCode,
    status: 1,
    page: pagination.value.page,
    pageSize: pagination.value.pageSize
  }
  const data = await getRoleList(dataParams)
  console.log(data.data);
  allRoles.value = data.data.list
  pagination.value.total = data.data.total

  // 角色列表加载完成后初始化已选角色
  initSelectedRoles()
}

const resetRoleList = () => {
  allRoles.value = []
}
// 计算属性
const dialogTitle = computed(() => {
  return props.title
})

// 过滤后的角色列表
const filteredRoles = computed(() => {
  let result = allRoles.value.filter(role => role.status === 1)
  return result
})

// 判断角色是否已选择
const isRoleSelected = (roleId) => {
  return !selectedRoles.value.some(r => r.id === roleId)
}

// 获取角色标签类型
const getRoleTagType = (role) => {
  const typeMap = {
    'admin': 'danger',
    'editor': 'success',
    'auditor': 'info'
  }
  return typeMap[role.code] || 'info'
}

// 监听 modelValue 变化
watch(() => props.modelValue, (newVal) => {
  console.log('modelValue 变化:', newVal, 'userId:', props.userId, 'existingRoleIds:', props.existingRoleIds)
  visible.value = newVal
  if (newVal) {
    // 打开对话框时先转换数据
    translate()
    // 角色列表加载后会在 loadRoleList 中调用 initSelectedRoles
  }
})

watch(visible, (newVal) => {
  emit('update:modelValue', newVal)
})


// 初始化已选角色
const initSelectedRoles = () => {
  console.log('initSelectedRoles - allRoles:', allRoles.value, 'roleIds:', roleIds.value)
  const existingRoles = allRoles.value.filter(role =>
    roleIds.value.includes(role.id)
  )
  console.log('initSelectedRoles - existingRoles:', existingRoles)
  selectedRoles.value = [...existingRoles]
}

// 搜索
const handleSearch = () => {
  pagination.value.page = 1
  loadRoleList()
}

// 重置搜索
const handleReset = () => {
  searchForm.value = {
    roleName: '',
    roleCode: ''
  }
  pagination.value.page = 1
}

// 表格选择变化
const handleSelectionChange = (selection) => {
  tableSelection.value = selection
  console.log('handleSelectionChange - selection:', selection)

  // 将表格选择的角色添加到已选角色列表(去重)
  selection.forEach(role => {
    if (!selectedRoles.value.some(r => r.id === role.id)) {
      selectedRoles.value.push(role)
    }
  })

  // 从已选角色中移除未被选择的
  const selectedIds = selection.map(r => r.id)
  selectedRoles.value = selectedRoles.value.filter(r => selectedIds.includes(r.id) || roleIds.value.includes(r.id))

  console.log('handleSelectionChange - selectedRoles:', selectedRoles.value)
}

// 添加角色
const handleAddRole = (role) => {
  if (!selectedRoles.value.some(r => r.id === role.id)) {
    selectedRoles.value.push(role)
  }
  // 清除表格选中状态
  if (roleTableRef.value) {
    roleTableRef.value.clearSelection()
  }
}

// 移除角色
const handleRemoveRole = (roleId) => {
  const index = selectedRoles.value.findIndex(r => r.id === roleId)
  if (index !== -1) {
    selectedRoles.value.splice(index, 1)
  }
}

// 分页变化
const handleSizeChange = (size) => {
  pagination.value.pageSize = size
}

const handleCurrentChange = (page) => {
  pagination.value.page = page
}

// 取消
const handleCancel = () => {
  visible.value = false
  emit('cancel')
}

// 确认
const handleConfirm = async () => {
  submitting.value = true
  try {
    const roleIds = selectedRoles.value.map(role => role.id)
    console.log('分配角色 - userId:', props.userId, 'roleIds:', roleIds, 'selectedRoles:', selectedRoles.value)
    emit('confirm', {
      userId: props.userId,
      roleIds: roleIds
    })
    visible.value = false
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  // 不需要在这里调用 translate,因为在对话框打开时会调用
})
</script>

<style scoped lang="scss">
.role-assign-dialog {
  :deep(.el-dialog__header) {
    padding: 20px 24px;
    border-bottom: 1px solid #ebeef5;
    margin: 0;

    .el-dialog__title {
      font-size: 18px;
      font-weight: 600;
      color: #303133;
    }

    .el-dialog__headerbtn {
      top: 20px;

      .el-dialog__close {
        font-size: 20px;
        color: #909399;

        &:hover {
          color: #409eff;
        }
      }
    }
  }

  :deep(.el-dialog__body) {
    padding: 24px;
  }

  :deep(.el-dialog__footer) {
    padding: 16px 24px;
    border-top: 1px solid #ebeef5;
    margin: 0;
  }
}

.dialog-content {
  display: flex;
  gap: 20px;
  height: 420px;
}

.left-section {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 12px;

  .search-area {
    padding: 12px;
    border-radius: 6px;

    .search-form {
      width: 100%;

      .search-inline-form {
        display: flex;
        align-items: center;
        margin: 0;

        :deep(.el-form-item) {
          margin-bottom: 0;
          margin-right: 12px;

          &:last-child {
            margin-right: 0;
          }

          .el-form-item__label {
            font-size: 13px;
            color: #606266;
            font-weight: 500;
          }

          .form-input {
            width: 140px;

            :deep(.el-input__wrapper) {
              padding: 4px 11px;
            }
          }
        }
      }
    }
  }

  .role-list-area {
    flex: 1;
    overflow: hidden;

    .role-table {
      :deep(.el-table__header) {
        th {
          background-color: #f5f7fa;
        }
      }

      .disabled-text {
        color: #c0c4cc;
        font-size: 13px;
      }
    }
  }

  .pagination-area {
    display: flex;
    justify-content: flex-end;
    align-items: center;
    padding: 8px 0;

    :deep(.el-pagination) {
      font-size: 13px;

      .el-pager li {
        min-width: 28px;
        height: 28px;
        line-height: 26px;
        margin: 0 2px;
        font-size: 12px;
      }

      .btn-prev,
      .btn-next {
        padding: 0 6px;
        min-width: 28px;
        height: 28px;
      }

      .el-pagination__jump {
        margin-left: 8px;

        .el-input__wrapper {
          padding: 0 6px;
          width: 50px;
        }

        .el-input__inner {
          text-align: center;
        }
      }

      .el-pagination__total,
      .el-pagination__jump {
        font-size: 12px;
        color: #606266;
      }
    }
  }
}

.right-section {
  width: 280px;
  display: flex;
  flex-direction: column;
  border-left: 1px solid #ebeef5;
  padding-left: 20px;

  .selected-header {
    display: flex;
    align-items: center;
    gap: 8px;
    padding-bottom: 16px;
    border-bottom: 1px solid #ebeef5;
    margin-bottom: 16px;

    .selected-title {
      font-size: 15px;
      font-weight: 600;
      color: #303133;
    }

    .selected-count {
      font-size: 14px;
      color: #909399;
    }
  }

  .selected-list {
    flex: 1;
    overflow-y: auto;

    &::-webkit-scrollbar {
      width: 6px;
    }

    &::-webkit-scrollbar-thumb {
      background: #dcdfe6;
      border-radius: 3px;

      &:hover {
        background: #c0c4cc;
      }
    }

    &::-webkit-scrollbar-track {
      background: transparent;
    }
  }

  .selected-items {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .selected-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 12px;
    background: #f5f7fa;
    border-radius: 6px;
    transition: all 0.2s;

    &:hover {
      background: #ecf5ff;
    }

    .role-info {
      display: flex;
      flex-direction: column;
      gap: 6px;
      flex: 1;

      .role-code {
        font-size: 12px;
        color: #909399;
      }
    }
  }
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;

  .el-button {
    padding: 10px 28px;
  }
}
</style>
