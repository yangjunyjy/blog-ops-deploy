<template>
  <div class="menu-list">
    <el-card shadow="never">
      <!-- 搜索栏 -->
      <el-form :model="searchForm" inline class="search-form" @submit.prevent>
        <div class="form-items">
          <el-form-item label="菜单名称">
            <el-input v-model="searchForm.title" placeholder="请输入菜单名称" clearable style="width: 200px" />
          </el-form-item>
          <el-form-item label="菜单类型">
            <el-select v-model="searchForm.type" placeholder="请选择" clearable style="width: 120px">
              <el-option label="全部" value="" />
              <el-option label="目录" value="directory" />
              <el-option label="菜单" value="menu" />
              <el-option label="按钮" value="button" />
            </el-select>
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
          <span class="action-title">菜单管理</span>
          <div class="action-buttons">
            <HaloButton type="primary" size="medium" :icon="Plus" content="新增" @click="handleCreate" />
            <HaloButton type="default" size="medium" :content="isExpanded ? '折叠' : '展开'" @click="handleExpandAll" />
          </div>
        </div>

        <el-table :data="displayedMenuList" v-loading="loading" row-key="id" :tree-props="{ children: 'children' }"
          ref="tableRef" :header-cell-style="{ backgroundColor: '#f5f7fa' }">
          <el-table-column prop="menu_name" label="菜单名称" min-width="200" />
          <el-table-column label="图标" width="100">
            <template #default="{ row }">
              <el-icon v-if="row.icon">
                <component :is="row.icon" />
              </el-icon>
            </template>
          </el-table-column>
          <el-table-column label="类型" width="100">
            <template #default="{ row }">
              <el-tag :type="getTypeTagType(row.menu_type)" size="small">
                {{ getTypeLabel(row.menu_type) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="path" label="路由地址" min-width="150" />
          <el-table-column prop="component" label="组件地址" min-width="150" />
          <el-table-column prop="menu_code" label="权限标识" min-width="150" />
          <el-table-column label="是否隐藏" width="100">
            <template #default="{ row }">
              <el-tag :type="!row.is_visible ? 'warning' : 'success'" size="small">
                {{ !row.is_visible ? '隐藏' : '显示' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="remark" label="菜单描述" min-width="150" />
          <el-table-column prop="sort" label="排序" width="80" align="center" />
          <el-table-column label="状态" width="80">
            <template #default="{ row }">
              <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
                {{ row.status === 1 ? '启用' : '禁用' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="100" fixed="right">
            <template #default="{ row }">
              <div style="display: flex; align-items: center; gap: 8px;">
                <el-button type="primary" link size="small" @click="handleEdit(row)" class="icon-btn">
                  <el-icon class="action-icon">
                    <Edit />
                  </el-icon>
                </el-button>
                <el-button type="danger" link size="small" @click="handleDelete(row.id)" class="icon-btn">
                  <el-icon class="action-icon">
                    <Delete />
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
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="菜单类型" prop="type">
          <el-radio-group v-model="form.type">
            <el-radio :value="1">目录</el-radio>
            <el-radio :value="2">菜单</el-radio>
            <el-radio :value="3">按钮</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="上级菜单" prop="parent_id">
          <el-tree-select v-model="form.parent_id" :data="menuTreeSelectData"
            :props="{ label: 'menu_name', value: 'id' }" clearable placeholder="请选择上级菜单" check-strictly />
        </el-form-item>
        <el-form-item label="菜单名称" prop="title">
          <el-input v-model="form.name" placeholder="请输入菜单名称" />
        </el-form-item>
        <el-form-item label="图标" prop="icon" v-if="form.type !== 3">
          <el-input v-model="form.icon" placeholder="请输入图标名称" />
        </el-form-item>
        <el-form-item label="路径" prop="path" v-if="form.type !== 3">
          <el-input v-model="form.path" placeholder="请输入路由地址" />
        </el-form-item>
        <el-form-item label="组件路径" prop="component" v-if="form.type === 2">
          <el-input v-model="form.component" placeholder="请输入组件路径" />
        </el-form-item>
        <el-form-item label="是否隐藏" prop="isVisible">
          <el-radio-group v-model="form.isVisible">
            <el-radio :value="true">否</el-radio>
            <el-radio :value="false">是</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="权限标识" prop="permission" v-if="form.type === 3">
          <el-input v-model="form.code" placeholder="请输入权限标识，如：system:user:create" />
        </el-form-item>
        <el-form-item label="菜单描述" prop="description">
          <el-input v-model="form.description" placeholder="请输入描述信息" style="height: 60px;" />
        </el-form-item>
        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="form.sort" :min="0" :max="999" />
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
  </div>
</template>

<script setup>
defineOptions({
  name: 'MenuList'
})

import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Edit, Delete } from '@element-plus/icons-vue'
import HaloButton from '../../layout/components/HaloButton.vue'
import { createMenu, updateMenu, deleteMenu as deleteMenuApi, getAllMenus } from '@/api/menu'

const loading = ref(false)
const dialogVisible = ref(false)
const submitting = ref(false)
const formRef = ref(null)
const tableRef = ref(null)
const isExpanded = ref(false)
const allMenus = ref([])

// 分页
const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const allMenuList = ref([]) // 存储所有顶级目录及其子菜单

const searchForm = reactive({
  title: '',
  type: ''
})

const dialogTitle = ref('新增菜单')
const form = reactive({
  id: null,
  parent_id: 0,
  name: '',
  icon: '',
  type: 'menu',
  path: '',
  component: '',
  code: '',
  sort: 0,
  status: 1,
  isVisible: true,
  description: ''
})

const rules = {
  parent_id: [{ required: true, message: '请选择上级菜单', trigger: 'change' }],
  name: [{ required: true, message: '请输入菜单名称', trigger: 'blur' }],
  type: [{ required: true, message: '请选择菜单类型', trigger: 'change' }],
  path: [{ required: true, message: '请选择菜单类型', trigger: 'blur' }]
}

// 构建菜单树选择数据（排除按钮）
const menuTreeSelectData = computed(() => {
  const buildTree = (parentId) => {
    if (!allMenus.value || !Array.isArray(allMenus.value)) {
      return []
    }
    return allMenus.value
      .filter(m => m.parent_id === parentId && m.type !== 'button' && m.type !== 2 && m.status === 1)
      .sort((a, b) => a.sort - b.sort)
      .map(m => ({
        ...m,
        children: buildTree(m.id)
      }))
  }
  return buildTree(0)
})

// 构建菜单树（用于表格展示）
const buildMenuTree = (parentId, allMenus) => {
  if (!allMenus || !Array.isArray(allMenus)) {
    return []
  }
  return allMenus
    .filter(m => m.parent_id === parentId)
    .sort((a, b) => a.sort - b.sort)
    .map(m => ({
      ...m,
      children: buildMenuTree(m.id, allMenus)
    }))
}

// 当前页显示的顶级目录及其子菜单
const displayedMenuList = computed(() => {
  const start = (pagination.page - 1) * pagination.pageSize
  const end = start + pagination.pageSize
  return allMenuList.value.slice(start, end)
})

const getTypeLabel = (type) => {
  const map = {
    directory: '目录',
    menu: '菜单',
    button: '按钮',
    1: '目录',
    2: '菜单',
    3: '按钮'
  }
  return map[type] || type
}

const getTypeTagType = (type) => {
  const map = {
    directory: 'warning',
    menu: 'primary',
    button: 'success',
    1: 'warning',
    2: 'primary',
    3: 'success'
  }
  return map[type] || 'info'
}

const loadList = async () => {
  // 获取所有菜单(扁平化结构)
  const data = await getAllMenus()
  console.log("获取的菜单数据为", data.data);

  // 保存所有菜单数据供递归使用
  allMenus.value = data.data

  allMenuList.value = allMenus.value
    .filter(m => m.parent_id === 0)
    .sort((a, b) => a.sort - b.sort)
    .map(m => ({
      ...m,
      children: buildMenuTree(m.id, allMenus.value)
    }))
  console.log("过滤后的菜单树为", allMenuList.value);
  // 分页统计的是顶级目录数量
  pagination.total = allMenuList.value.length
  pagination.page = 1

}

const handleSearch = () => {
  pagination.page = 1
  if (searchForm.title || searchForm.type) {
    const filtered = allMenus.value.filter(m => {
      const matchTitle = !searchForm.title || m.menu_name?.includes(searchForm.title)
      const matchType = !searchForm.type || m.menu_type === searchForm.type
      return matchTitle && matchType
    })

    // 获取匹配的顶级目录
    const matchedTopMenus = allMenus.value.filter(m => {
      if (m.parent_id !== 0) return false
      return filterTree(filtered, m.id, allMenus.value)
    })

    allMenuList.value = matchedTopMenus
      .sort((a, b) => a.sort - b.sort)
      .map(m => ({
        ...m,
        children: buildMenuTree(m.id, allMenus.value)
      }))

    pagination.total = allMenuList.value.length
  } else {
    loadList()
  }
}

const filterTree = (filteredMenus, id, allMenus) => {
  const menu = filteredMenus.find(m => m.id === id)
  if (menu) return true
  // 检查子菜单
  const children = allMenus.filter(m => m.parent_id === id)
  return children.some(child => filterTree(filteredMenus, child.id, allMenus))
}

const handleReset = () => {
  searchForm.title = ''
  searchForm.type = ''
  pagination.page = 1
  loadList()
}

const handleCreate = () => {
  dialogTitle.value = '新增菜单'
  Object.assign(form, {
    id: null,
    parent_id: 0,
    name: '',
    icon: '',
    type: 2,
    path: '',
    component: '',
    code: '',
    sort: 0,
    status: 1,
    isVisible: true,
    description: ''
  })
  dialogVisible.value = true
}

const handleEdit = (row) => {
  dialogTitle.value = '编辑菜单'
  Object.assign(form, {
    id: row.id,
    parent_id: row.parent_id,
    path: row.path,
    status: row.status,
    sort: row.sort,
    isVisible: row.is_visible,
    component: row.component,
    description: row.remark,
    type: row.menu_type,
    name: row.menu_name,
    code: row.menu_code,
    icon: row.icon
  })
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitting.value = true
    try {
      // 构建提交数据，确保布尔值正确
      const submitData = {
        ...form,
        isVisible: form.isVisible === true || form.isVisible === 'true' ? true : false
      }

      if (form.id) {
        const data = await updateMenu(submitData)
        if (data.code == 200) {
          ElMessage.success('更新成功')
        }
      } else {
        const data = await createMenu(submitData)
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
    await ElMessageBox.confirm('确定删除该菜单吗？如果有子菜单，子菜单也会被删除。', '提示', { type: 'warning' })

    // 调用删除API
    const data = await deleteMenuApi(id)

    if (data.code == 200) {
      ElMessage.success('删除成功')
      loadList()
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const handleExpandAll = () => {
  isExpanded.value = !isExpanded.value

  // 递归展开/折叠所有行
  const toggleExpand = (rows) => {
    if (!rows) return
    rows.forEach(row => {
      tableRef.value.toggleRowExpansion(row, isExpanded.value)
      if (row.children && row.children.length > 0) {
        toggleExpand(row.children)
      }
    })
  }

  toggleExpand(allMenuList.value)
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
.menu-list {
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
}

.icon-btn {
  margin: 0 4px;

  &:hover .action-icon {
    background-color: rgba(64, 158, 255, 0.1);
  }
}
</style>
