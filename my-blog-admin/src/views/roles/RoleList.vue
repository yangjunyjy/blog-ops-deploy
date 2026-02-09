<template>
  <div class="role-list">
    <el-card shadow="never">
      <!-- 搜索栏 -->
      <el-form :model="searchForm" inline class="search-form" @submit.prevent>
        <div class="form-items">
          <el-form-item label="角色名称">
            <el-input v-model="searchForm.name" placeholder="角色名称" clearable style="width: 150px" />
          </el-form-item>
          <el-form-item label="角色编码">
            <el-input v-model="searchForm.code" placeholder="角色编码" clearable style="width: 150px" />
          </el-form-item>
          <el-select v-model="searchForm.status" placeholder="状态" style="width: 80px" clearable :teleported="false">
            <el-option label="在用" :value="1" />
            <el-option label="禁用" :value="0" />
          </el-select>
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
          <span class="action-title">角色管理</span>
          <div class="action-buttons">
            <HaloButton type="primary" size="medium" :icon="Plus" content="新增" @click="handleCreate" />
          </div>
        </div>

        <el-table :data="roleList" v-loading="loading" :header-cell-style="{ backgroundColor: '#f5f7fa' }">
          <el-table-column prop="id" label="ID" width="70" align="center" />
          <el-table-column prop="roleName" label="角色名称" min-width="120" />
          <el-table-column prop="roleCode" label="角色编码" min-width="120" />
          <el-table-column prop="roleDesc" label="描述" min-width="200" />
          <el-table-column label="状态" width="80">
            <template #default="{ row }">
              <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
                {{ row.status === 1 ? '启用' : '禁用' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="140" fixed="right">
            <template #default="{ row }">
              <div style="display: flex; align-items: center; gap: 8px;">
                <el-button type="primary" link size="small" @click="handleEdit(row)" class="icon-btn">
                  <el-icon class="action-icon">
                    <img src="/svg/edit.svg" alt="edit" style="width: 16px; height: 16px;" />
                  </el-icon>
                </el-button>
                <el-button type="success" link size="small" @click="handleAssign(row)" class="icon-btn">
                  <el-icon class="action-icon">
                    <img src="/svg/user.svg" alt="user" style="width: 16px; height: 16px;" />
                  </el-icon>
                </el-button>
                <el-button type="danger" link size="small" @click="handleDelete(row.id)" class="icon-btn">
                  <el-icon class="action-icon delete-icon">
                    <img src="/svg/delete.svg" alt="delete" style="width: 16px; height: 16px; color: #f56c6c;"
                      class="svg-red" />
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
        <el-form-item label="角色名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入角色名称" />
        </el-form-item>
        <el-form-item label="角色编码" prop="code">
          <el-input v-model="form.code" placeholder="请输入角色编码" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="form.description" type="textarea" :rows="3" placeholder="请输入描述" />
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

    <!-- 分配权限对话框 -->
    <el-dialog v-model="assignDialogVisible" title="分配权限" width="600px" @open="GetRoleMenus(currentRole.id)"
      @close="DestroyLastRoleData">
      <el-form :model="assignForm" ref="assignFormRef" label-width="100px">
        <el-form-item label="角色名称">
          <span>{{ currentRole?.roleName }}</span>
        </el-form-item>
        <el-form-item label="权限">
          <el-tree ref="menuTreeRef" :data="menuTree" :props="{ label: 'title', children: 'children' }" node-key="id"
            show-checkbox :default-checked-keys="checkedMenuIds" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="assignDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleAssignSubmit" :loading="assigning">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
defineOptions({
  name: 'RoleList'
})

import { ref, reactive, onMounted, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import HaloButton from '../../layout/components/HaloButton.vue'
import { useRbacStore } from '../../store/rbac'
import { assignRoleMenus, getAllMenus, GetMenusByRoleID, getRoleList, createRole, updateRole, deleteRole } from '@/api'


const loading = ref(false)
const dialogVisible = ref(false)
const assignDialogVisible = ref(false)
const submitting = ref(false)
const assigning = ref(false)
const formRef = ref(null)
const assignFormRef = ref(null)
const menuTree = ref([])
const menuTreeRef = ref(null)
const roleList = ref([])
const currentRole = ref(null)
const checkedMenuIds = ref([])
const rbacStore = useRbacStore()

// 分页
const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const allRoleList = ref([]) // 存储所有角色数据

const searchForm = reactive({
  name: '',
  code: '',
  status: 1
})

const dialogTitle = ref('新增角色')
const form = reactive({
  id: null,
  name: '',
  code: '',
  description: '',
  status: 1
})

const assignForm = reactive({
  menuIds: []
})

const rules = {
  name: [{ required: true, message: '请输入角色名称', trigger: 'blur' }],
  code: [{ required: true, message: '请输入角色编码', trigger: 'blur' }]
}


const loadList = async () => {
  const params = {
    roleName: searchForm.name,
    roleCode: searchForm.code,
    status: searchForm.status,
    page: pagination.page,
    pageSize: pagination.pageSize
  }
  const data = await getRoleList(params)
  console.log(data.data);
  pagination.total = data.data.total
  roleList.value = data.data.list
}

const handleSearch = () => {
  pagination.page = 1
  if (searchForm.name) {
    allRoleList.value = rbacStore.roles.filter(r =>
      r.name.includes(searchForm.name)
    )
    pagination.total = allRoleList.value.length
  } else {
    loadList()
  }
}

const handleReset = () => {
  searchForm.name = ''
  pagination.page = 1
  searchForm.code = ''
  searchForm.status = 1
  loadList()
}

const handleCreate = () => {
  dialogTitle.value = '新增角色'
  Object.assign(form, {
    id: null,
    name: '',
    code: '',
    description: '',
    status: 1
  })
  dialogVisible.value = true
}

const handleEdit = (row) => {
  dialogTitle.value = '编辑角色'
  Object.assign(form, {
    id: row.id,
    name: row.roleName,
    code: row.roleCode,
    description: row.roleDesc,
    status: row.status
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
        const params = {
          id: form.id,
          role_name: form.name,
          role_code: form.code,
          status: form.status,
          role_desc: form.description
        }
        const data = await updateRole(params)
        console.log("更新角色的响应", data);
        if (data.code == 200) {
          ElMessage.success('更新成功')
        }
      } else {
        const params = {
          role_code: form.code,
          role_name: form.name,
          role_desc: form.description,
          status: form.status
        }
        const data = await createRole(params)
        console.log("创建角色的响应数据", data);

        if (data.code == 200) {
          ElMessage.success('创建成功')
        }
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
    await ElMessageBox.confirm('确定删除该角色吗？', '提示', { type: 'warning' })
    const data = await deleteRole(id)
    if (!data) return
    ElMessage.success('删除成功')
    loadList()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const handleAssign = (row) => {
  currentRole.value = row
  checkedMenuIds.value = []  // 先清空，等待对话框打开时重新获取
  assignDialogVisible.value = true
}


const handleAssignSubmit = async () => {
  assigning.value = true
  try {
    // 获取所有完全选中的节点（包括父节点自动选中的子节点）
    const checkedKeys = menuTreeRef.value.getCheckedKeys()
    // 获取半选中的节点（父节点下只有部分子节点被选中）
    const halfCheckedKeys = menuTreeRef.value.getHalfCheckedKeys()
    // 合并：完全选中的节点 + 半选中的父节点
    const allMenuIds = [...checkedKeys, ...halfCheckedKeys]

    console.log("完全选中的节点id（包括自动选中的子节点）", checkedKeys);
    console.log("半选中的父节点id", halfCheckedKeys);
    console.log("最终提交的菜单id", allMenuIds);

    const roleId = currentRole.value.id
    console.log("当前角色id", roleId);
    const params = {
      role_id: roleId,
      menu_ids: allMenuIds
    }
    const data = await assignRoleMenus(params)
    console.log("分配响应数据", data);

    if (data.code == 200) {
      ElMessage.success('分配权限成功')
    }
    assignDialogVisible.value = false
  } catch (error) {
    ElMessage.error('分配权限失败')
  } finally {
    assigning.value = false
  }
}

const handleSizeChange = (size) => {
  pagination.pageSize = size
  pagination.page = 1
}

const handleCurrentChange = (page) => {
  pagination.page = page
}

const GetRoleMenus = async (id) => {
  const dataRole = await GetMenusByRoleID(id)
  console.log("角色的菜单为", dataRole.data);
  const dataAll = await getAllMenus()
  console.log("所有菜单数据为", dataAll.data)
  // 转化菜单数据为树形结构
  menuTree.value = buildMenusTree(dataAll.data)
  // 转化后的树形结构为
  console.log("转化后的树形结构", menuTree.value)

  // 只保留叶子节点（没有子菜单的节点）在选中列表中
  // 这样父菜单会根据子菜单的选中状态自动显示为半选中或全选中
  const isLeafNode = (menuId, allMenus) => {
    return !allMenus.some(m => m.parent_id === menuId)
  }

  checkedMenuIds.value = dataRole.data
    .filter(menu => isLeafNode(menu.id, dataAll.data))
    .map(menu => menu.id)

  console.log("过滤后只保留叶子节点的菜单id数组", checkedMenuIds.value);

  // 等待树形结构渲染完后设置选中状态
  nextTick(() => {
    if (menuTreeRef.value) {
      menuTreeRef.value.setCheckedKeys(checkedMenuIds.value)
    }
  })
}

// 把扁平化的菜单数组转化成菜单树
const buildMenusTree = (menus, parent_id = 0) => {
  return menus.filter(menu => menu.parent_id == parent_id && menu.status == 1).map(menu => ({
    ...menu,
    title: menu.menu_name,
    children: buildMenusTree(menus, menu.id)
  }))
}

const DestroyLastRoleData = () => {
  if (menuTreeRef.value) {
    menuTreeRef.value.setCheckedKeys([])
  }
  menuTree.value = []
  checkedMenuIds.value = []
  console.log("关闭时清空上个角色数据", menuTree.value, checkedMenuIds.value);
}

onMounted(() => {
  loadList()
})
</script>

<style scoped lang="scss">
.role-list {
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
    gap: 0;
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

  img {
    vertical-align: middle;
  }

  &.delete-icon img {
    filter: brightness(0) saturate(100%) invert(22%) sepia(87%) saturate(3827%) hue-rotate(352deg) brightness(92%) contrast(85%);
  }
}

.icon-btn {
  margin: 0 4px;

  &:hover .action-icon {

    transform: scale(1.1);

    background-color: rgba(64, 158, 255, 0.1);

  }

}
</style>