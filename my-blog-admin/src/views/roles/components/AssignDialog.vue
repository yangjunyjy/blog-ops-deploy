<template>
    <!-- 分配权限对话框 -->
    <el-dialog v-model="assignDialogVisible" :title="`分配权限 - ${currentRole?.roleName || ''}`" width="780px"
        :close-on-click-modal="false" @closed="handleDialogClosed">
        <!-- 权限分配面板 -->
        <div class="permission-assign-panel">
            <!-- 左侧权限树 -->
            <div class="permission-tree-container">
                <div class="tree-header">
                    <div class="header-left">
                        <span class="title">菜单权限</span>
                        <el-tag v-if="checkedCount > 0" type="info" size="small">
                            已选 {{ checkedCount }} 项
                        </el-tag>
                    </div>
                    <div class="header-right">
                        <el-input v-model="treeFilterText" placeholder="搜索菜单..." clearable size="small"
                            style="width: 160px" @input="filterTree">
                            <template #prefix>
                                <el-icon>
                                    <Search />
                                </el-icon>
                            </template>
                        </el-input>
                        <el-tooltip content="全选">
                            <el-button type="text" size="small" @click="selectAll" :disabled="!menuTreeData.length">
                                全选
                            </el-button>
                        </el-tooltip>
                        <el-tooltip content="清空">
                            <el-button type="text" size="small" @click="clearAll" :disabled="!menuTreeData.length">
                                清空
                            </el-button>
                        </el-tooltip>
                    </div>
                </div>

                <div class="tree-wrapper">
                    <el-scrollbar height="380px">
                        <el-tree ref="menuTreeRef" :data="filteredMenuTree" :props="treeProps" node-key="id"
                            show-checkbox :default-checked-keys="defaultCheckedKeys"
                            :filter-node-method="filterTreeNode" highlight-current :expand-on-click-node="false"
                            @check="handleTreeCheck">
                            <template #default="{ node, data }">
                                <div class="tree-node-content">
                                    <div class="node-left">
                                        <el-icon v-if="data.icon" class="node-icon">
                                            <component :is="data.icon" />
                                        </el-icon>
                                        <span class="node-label">{{ data.title }}</span>
                                        <el-tag v-if="data.type === 2" size="mini" type="danger" class="perm-tag">
                                            按钮
                                        </el-tag>
                                        <el-tag v-if="data.hidden" size="mini" type="info" class="perm-tag">
                                            隐藏
                                        </el-tag>
                                    </div>
                                    <div v-if="data.perms" class="node-right">
                                        <el-tooltip :content="data.perms">
                                            <el-tag size="mini" class="perm-code">
                                                {{ truncatePerm(data.perms) }}
                                            </el-tag>
                                        </el-tooltip>
                                    </div>
                                </div>
                            </template>
                        </el-tree>
                    </el-scrollbar>
                </div>
            </div>

            <!-- 右侧权限说明 -->
            <div class="permission-info">
                <div class="info-header">
                    <el-icon class="info-icon">
                        <InfoFilled />
                    </el-icon>
                    <span class="info-title">权限说明</span>
                </div>

                <div class="info-content">
                    <div class="role-info">
                        <div class="info-item">
                            <span class="label">角色名称：</span>
                            <span class="value">{{ currentRole?.roleName || '-' }}</span>
                        </div>
                        <div class="info-item">
                            <span class="label">角色编码：</span>
                            <span class="value code">{{ currentRole?.roleCode || '-' }}</span>
                        </div>
                        <div class="info-item">
                            <span class="label">角色描述：</span>
                            <span class="value">{{ currentRole?.roleDesc || '暂无描述' }}</span>
                        </div>
                    </div>

                    <el-divider />

                    <div class="permission-tips">
                        <h4 class="tips-title">权限分配提示</h4>
                        <ul class="tips-list">
                            <li class="tip-item">
                                <el-icon class="tip-icon">
                                    <CircleCheck />
                                </el-icon>
                                <span>选中菜单节点，用户可访问该页面</span>
                            </li>
                            <li class="tip-item">
                                <el-icon class="tip-icon">
                                    <CircleCheck />
                                </el-icon>
                                <span>"按钮"权限控制页面内操作按钮的显示</span>
                            </li>
                            <li class="tip-item">
                                <el-icon class="tip-icon">
                                    <CircleCheck />
                                </el-icon>
                                <span>"隐藏"标签表示菜单在前端不显示</span>
                            </li>
                            <li class="tip-item">
                                <el-icon class="tip-icon">
                                    <CircleCheck />
                                </el-icon>
                                <span>鼠标悬停可查看完整的权限标识</span>
                            </li>
                        </ul>
                    </div>

                    <div class="selected-preview" v-if="checkedCount > 0">
                        <h4 class="preview-title">
                            已选权限预览
                            <span class="preview-count">({{ checkedCount }})</span>
                        </h4>
                        <div class="preview-tags">
                            <el-tag v-for="menu in selectedMenus" :key="menu.id" size="small" closable
                                @close="removeMenu(menu.id)" class="preview-tag">
                                {{ menu.title }}
                                <span v-if="menu.type === 2" class="tag-type">(按钮)</span>
                            </el-tag>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <template #footer>
            <div class="dialog-footer">
                <div class="footer-left">
                    <el-button type="text" :icon="Refresh" @click="resetSelection"
                        :disabled="!menuTreeData.length || assigning">
                        重置
                    </el-button>
                </div>
                <div class="footer-right">
                    <el-button @click="assignDialogVisible = false" :disabled="assigning">
                        取消
                    </el-button>
                    <el-button type="primary" @click="handleAssignSubmit" :loading="assigning"
                        :disabled="checkedCount === 0">
                        确认分配
                    </el-button>
                </div>
            </div>
        </template>
    </el-dialog>
</template>

<script setup>
import { ref, computed, watch, nextTick } from 'vue'
import {
    Search,
    InfoFilled,
    CircleCheck,
    Refresh,
    Menu as MenuIcon,
    Document,
    Setting,
    User,
    Tools,
    DataAnalysis
} from '@element-plus/icons-vue'

const props = defineProps({
    visible: Boolean,
    currentRole: Object,
    menuTreeData: Array,
    defaultCheckedKeys: Array
})

const emit = defineEmits(['update:visible', 'assign-submit'])

const assignDialogVisible = computed({
    get: () => props.visible,
    set: (value) => emit('update:visible', value)
})

const menuTreeRef = ref(null)
const treeFilterText = ref('')
const assigning = ref(false)

// 树形配置
const treeProps = {
    label: 'title',
    children: 'children',
    disabled: (data) => data.disabled || false
}

// 图标映射（根据实际情况调整）
const iconMap = {
    1: MenuIcon,
    2: Document,
    3: Setting,
    4: User,
    5: Tools,
    6: DataAnalysis
}

// 处理菜单图标
const processMenuTree = (tree) => {
    return tree.map(item => {
        const processed = {
            ...item,
            icon: iconMap[item.type] || MenuIcon
        }
        if (item.children && item.children.length) {
            processed.children = processMenuTree(item.children)
        }
        return processed
    })
}

// 处理后的菜单树
const processedMenuTree = computed(() => {
    return processMenuTree(props.menuTreeData)
})

// 过滤后的菜单树
const filteredMenuTree = computed(() => {
    if (!treeFilterText.value) return processedMenuTree.value
    return filterTreeData(processedMenuTree.value, treeFilterText.value)
})

// 选中的菜单数量
const checkedCount = computed(() => {
    if (!menuTreeRef.value) return 0
    const checkedNodes = menuTreeRef.value.getCheckedNodes(false, true)
    return checkedNodes.length
})

// 选中的菜单列表（用于预览）
const selectedMenus = computed(() => {
    if (!menuTreeRef.value) return []
    return menuTreeRef.value.getCheckedNodes(false, true).slice(0, 10) // 最多显示10个
})

// 过滤树数据
const filterTreeData = (data, filterText) => {
    return data.filter(node => {
        if (node.title?.toLowerCase().includes(filterText.toLowerCase())) {
            return true
        }
        if (node.children && node.children.length) {
            node.children = filterTreeData(node.children, filterText)
            return node.children.length > 0
        }
        return false
    })
}

// 过滤树节点
const filterTreeNode = (value, data) => {
    if (!value) return true
    return data.title?.toLowerCase().includes(value.toLowerCase())
}

// 树过滤
const filterTree = () => {
    nextTick(() => {
        if (menuTreeRef.value) {
            menuTreeRef.value.filter(treeFilterText.value)
        }
    })
}

// 全选
const selectAll = () => {
    if (menuTreeRef.value) {
        const allKeys = getAllKeys(processedMenuTree.value)
        menuTreeRef.value.setCheckedKeys(allKeys)
    }
}

// 清空
const clearAll = () => {
    if (menuTreeRef.value) {
        menuTreeRef.value.setCheckedKeys([])
    }
}

// 重置选择
const resetSelection = () => {
    if (menuTreeRef.value) {
        menuTreeRef.value.setCheckedKeys(props.defaultCheckedKeys || [])
        treeFilterText.value = ''
    }
}

// 获取所有节点key
const getAllKeys = (tree) => {
    let keys = []
    tree.forEach(node => {
        keys.push(node.id)
        if (node.children && node.children.length) {
            keys = keys.concat(getAllKeys(node.children))
        }
    })
    return keys
}

// 处理树节点选中
const handleTreeCheck = (checkedData, checkedInfo) => {
    // 可以在这里处理父子节点关联逻辑
    console.log('选中数据:', checkedData, checkedInfo)
}

// 移除菜单
const removeMenu = (menuId) => {
    if (menuTreeRef.value) {
        menuTreeRef.value.setChecked(menuId, false, false)
    }
}

// 截断权限标识
const truncatePerm = (perm) => {
    if (!perm) return ''
    return perm.length > 15 ? perm.substring(0, 15) + '...' : perm
}

// 处理对话框关闭
const handleDialogClosed = () => {
    treeFilterText.value = ''
    assigning.value = false
}

// 处理分配提交
const handleAssignSubmit = async () => {
    if (!menuTreeRef.value) return

    const checkedKeys = menuTreeRef.value.getCheckedKeys()
    const halfCheckedKeys = menuTreeRef.value.getHalfCheckedKeys()

    // 如果需要半选中的父节点key（例如只选中了子节点，父节点需要半选中状态）
    const allCheckedKeys = [...checkedKeys, ...halfCheckedKeys]

    assigning.value = true

    try {
        await emit('assign-submit', {
            roleId: props.currentRole?.id,
            menuIds: allCheckedKeys,
            checkedKeys,
            halfCheckedKeys
        })
        assignDialogVisible.value = false
    } catch (error) {
        console.error('分配权限失败:', error)
    } finally {
        assigning.value = false
    }
}

// 监听对话框打开
watch(() => props.visible, (visible) => {
    if (visible) {
        nextTick(() => {
            if (menuTreeRef.value) {
                // 确保树组件渲染后设置默认选中
                menuTreeRef.value.setCheckedKeys(props.defaultCheckedKeys || [])
            }
        })
    }
})
</script>

<style scoped lang="scss">
.permission-assign-panel {
    display: flex;
    gap: 20px;
    min-height: 460px;
}

.permission-tree-container {
    flex: 1;
    border: 1px solid #ebeef5;
    border-radius: 8px;
    overflow: hidden;
    background: #fafbfc;

    .tree-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 12px 16px;
        background: linear-gradient(135deg, #f5f7fa 0%, #e4e7ed 100%);
        border-bottom: 1px solid #ebeef5;

        .header-left {
            display: flex;
            align-items: center;
            gap: 12px;

            .title {
                font-size: 14px;
                font-weight: 600;
                color: #303133;
            }
        }

        .header-right {
            display: flex;
            align-items: center;
            gap: 8px;
        }
    }

    .tree-wrapper {
        padding: 12px;

        :deep(.el-tree) {
            background: transparent;

            .el-tree-node {
                margin: 4px 0;

                &:focus>.el-tree-node__content {
                    background-color: #ecf5ff;
                }

                .el-tree-node__content {
                    height: 36px;
                    border-radius: 4px;
                    transition: all 0.2s;

                    &:hover {
                        background-color: #f5f7fa;
                    }

                    .el-checkbox {
                        margin-right: 8px;
                    }
                }

                .el-tree-node__children {
                    padding-left: 16px;
                }
            }

            .tree-node-content {
                display: flex;
                align-items: center;
                justify-content: space-between;
                width: 100%;

                .node-left {
                    display: flex;
                    align-items: center;
                    gap: 6px;

                    .node-icon {
                        font-size: 14px;
                        color: #909399;
                    }

                    .node-label {
                        font-size: 13px;
                        color: #606266;
                    }

                    .perm-tag {
                        height: 18px;
                        line-height: 16px;
                        font-size: 10px;
                        padding: 0 4px;
                    }
                }

                .node-right {
                    .perm-code {
                        background: rgba(64, 158, 255, 0.1);
                        border-color: rgba(64, 158, 255, 0.2);
                        color: #409eff;
                        font-family: 'Monaco', 'Consolas', monospace;
                        font-size: 10px;
                        padding: 0 6px;
                        height: 20px;
                        line-height: 18px;
                    }
                }
            }
        }
    }
}

.permission-info {
    width: 280px;
    border: 1px solid #ebeef5;
    border-radius: 8px;
    overflow: hidden;
    background: #fafbfc;

    .info-header {
        display: flex;
        align-items: center;
        gap: 8px;
        padding: 12px 16px;
        background: linear-gradient(135deg, #f0f9ff 0%, #e6f7ff 100%);
        border-bottom: 1px solid #ebeef5;

        .info-icon {
            color: #409eff;
            font-size: 16px;
        }

        .info-title {
            font-size: 14px;
            font-weight: 600;
            color: #303133;
        }
    }

    .info-content {
        padding: 16px;

        .role-info {
            .info-item {
                margin-bottom: 10px;
                font-size: 13px;

                .label {
                    color: #909399;
                    min-width: 60px;
                    display: inline-block;
                }

                .value {
                    color: #606266;
                    word-break: break-all;

                    &.code {
                        font-family: 'Monaco', 'Consolas', monospace;
                        background: #f0f2f5;
                        padding: 2px 6px;
                        border-radius: 3px;
                    }
                }
            }
        }

        .permission-tips {
            margin-top: 20px;

            .tips-title {
                font-size: 13px;
                color: #303133;
                margin-bottom: 12px;
                font-weight: 600;
            }

            .tips-list {
                list-style: none;
                padding: 0;
                margin: 0;

                .tip-item {
                    display: flex;
                    align-items: flex-start;
                    gap: 8px;
                    margin-bottom: 10px;
                    font-size: 12px;
                    color: #606266;
                    line-height: 1.4;

                    .tip-icon {
                        color: #67c23a;
                        font-size: 12px;
                        margin-top: 2px;
                        flex-shrink: 0;
                    }
                }
            }
        }

        .selected-preview {
            margin-top: 20px;

            .preview-title {
                font-size: 13px;
                color: #303133;
                margin-bottom: 12px;
                font-weight: 600;

                .preview-count {
                    color: #909399;
                    font-weight: normal;
                    font-size: 12px;
                }
            }

            .preview-tags {
                display: flex;
                flex-wrap: wrap;
                gap: 6px;
                max-height: 120px;
                overflow-y: auto;
                padding-right: 4px;

                .preview-tag {
                    font-size: 11px;
                    transition: all 0.2s;

                    .tag-type {
                        color: #f56c6c;
                        font-size: 10px;
                    }

                    &:hover {
                        transform: translateY(-1px);
                        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
                    }
                }
            }
        }
    }
}

.dialog-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding-top: 16px;
    border-top: 1px solid #ebeef5;

    .footer-left {
        .el-button {
            color: #909399;

            &:hover {
                color: #409eff;
            }
        }
    }

    .footer-right {
        display: flex;
        gap: 12px;
    }
}

:deep(.el-dialog__header) {
    padding: 16px 20px;
    border-bottom: 1px solid #ebeef5;
    margin-right: 0;

    .el-dialog__title {
        font-size: 16px;
        font-weight: 600;
        color: #303133;
    }
}

:deep(.el-dialog__body) {
    padding: 20px;
}

:deep(.el-divider) {
    margin: 16px 0;
}

// 滚动条样式
:deep(.el-scrollbar__bar) {
    opacity: 0.6;
}
</style>