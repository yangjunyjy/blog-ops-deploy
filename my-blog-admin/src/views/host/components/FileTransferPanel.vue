<template>
    <el-drawer v-model="visible" title="文件传输" direction="rtl" size="33%" :close-on-click-modal="false"
        @close="handleClose" class="file-transfer-drawer">
        <div class="file-transfer-container">
            <!-- 顶部路径栏 -->
            <div class="header-section">
                <div class="path-input-wrapper">
                    <el-icon class="path-icon">
                        <Folder />
                    </el-icon>
                    <el-breadcrumb separator="/" class="path-breadcrumb">
                        <el-breadcrumb-item @click="navigateTo('/')">
                            <span>根目录</span>
                        </el-breadcrumb-item>
                        <el-breadcrumb-item v-for="(path, index) in pathParts" :key="index" @click="navigateTo(path)">
                            {{ path.split('/').filter(Boolean).pop() }}
                        </el-breadcrumb-item>
                    </el-breadcrumb>
                </div>
            </div>

            <!-- 操作工具栏 -->
            <div class="action-bar">
                <el-button-group class="btn-group-main">
                    <el-button size="small" :icon="Upload" type="primary" @click="handleUploadClick">
                        上传
                    </el-button>
                    <el-button size="small" :icon="Download" @click="handleDownloadClick" :disabled="!selectedFile">
                        下载
                    </el-button>
                </el-button-group>

                <el-divider direction="vertical" />

                <el-button-group class="btn-group-secondary">
                    <el-button size="small" :icon="FolderAdd" @click="handleMkdirClick">
                        新建
                    </el-button>
                    <el-button size="small" :icon="Delete" @click="handleDeleteClick" :disabled="!selectedFile"
                        type="danger">
                        删除
                    </el-button>
                    <el-button size="small" :icon="Refresh" @click="handleRefresh">
                        刷新
                    </el-button>
                </el-button-group>
            </div>

            <!-- 文件列表 -->
            <div class="file-list-wrapper">
                <el-table :data="files" @row-dblclick="handleRowDblClick" @row-click="handleRowClick"
                    highlight-current-row height="100%" size="small" stripe>
                    <el-table-column width="40" align="center">
                        <template #default="{ row }">
                            <el-icon v-if="row.isDir" color="#409eff" :size="16">
                                <Folder />
                            </el-icon>
                            <el-icon v-else color="#909399" :size="16">
                                <Document />
                            </el-icon>
                        </template>
                    </el-table-column>
                    <el-table-column label="文件名" min-width="80" show-overflow-tooltip>
                        <template #default="{ row }">
                            <span :class="{ 'parent-dir': row.name === '..' }">
                                {{ row.name === '..' ? '...' : row.name }}
                            </span>
                        </template>
                    </el-table-column>
                    <el-table-column label="大小" width="80" align="right">
                        <template #default="{ row }">
                            <span v-if="!row.isDir" class="file-size">{{ formatSize(row.size) }}</span>
                        </template>
                    </el-table-column>
                    <el-table-column label="修改时间" width="130" align="center">
                        <template #default="{ row }">
                            {{ formatTime(row.modTime) }}
                        </template>
                    </el-table-column>
                    <el-table-column prop="kind" label="文件类型" min-width="50" show-overflow-tooltip />
                </el-table>
            </div>

            <!-- 传输进度 - 固定在底部 -->
            <div class="transfer-section" v-show="transfers.length > 0">
                <div class="transfer-header">
                    <span class="transfer-title">
                        <el-icon>
                            <Upload />
                        </el-icon>
                        传输队列 ({{ transfers.length }})
                    </span>
                    <el-button link size="small" @click="clearCompletedTransfers" v-if="hasCompletedTransfers">
                        清除已完成
                    </el-button>
                </div>
                <div class="transfer-list">
                    <div v-for="transfer in transfers" :key="transfer.id" class="transfer-item">
                        <div class="transfer-info">
                            <el-icon class="transfer-icon">
                                <Upload v-if="transfer.type === 'upload'" />
                                <Download v-else />
                            </el-icon>
                            <div class="transfer-name">{{ transfer.name }}</div>
                        </div>
                        <div class="transfer-status">
                            <el-progress :percentage="Math.round(transfer.progress)"
                                :status="getProgressStatus(transfer.status)" :stroke-width="6" :show-text="false"
                                class="transfer-progress" />
                            <span class="transfer-percent">{{ Math.round(transfer.progress) }}%</span>
                            <el-button v-if="transfer.status === 'uploading' || transfer.status === 'downloading'" link
                                type="danger" size="small" @click="handleCancelTransfer(transfer.id)">
                                取消
                            </el-button>
                        </div>
                    </div>
                </div>
            </div>

            <!-- 隐藏的上传input -->
            <input ref="fileInputRef" type="file" style="display: none" multiple @change="handleFileSelect" />

            <!-- 对话框 -->
            <el-dialog v-model="renameDialogVisible" title="重命名" width="400px">
                <el-input v-model="newFileName" placeholder="请输入新文件名" />
                <template #footer>
                    <el-button @click="renameDialogVisible = false">取消</el-button>
                    <el-button type="primary" @click="handleRenameConfirm">确定</el-button>
                </template>
            </el-dialog>

            <el-dialog v-model="mkdirDialogVisible" title="新建目录" width="400px">
                <el-input v-model="newDirName" placeholder="请输入目录名" />
                <template #footer>
                    <el-button @click="mkdirDialogVisible = false">取消</el-button>
                    <el-button type="primary" @click="handleMkdirConfirm">确定</el-button>
                </template>
            </el-dialog>
        </div>
    </el-drawer>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
    Upload, Download, Delete, Refresh, Folder, Document,
    FolderAdd, HomeFilled, Edit
} from '@element-plus/icons-vue'
import { listFiles, uploadFile, downloadFile, deleteFile, createDir, renameFile } from '@/api/ssh_file'

defineOptions({
    name: 'FileTransferPanel'
})

const props = defineProps({
    visible: {
        type: Boolean,
        default: false
    },
    sessionId: {
        type: String,
        required: true
    }
})

const emit = defineEmits(['update:visible'])

// 状态
const visible = computed({
    get: () => props.visible,
    set: (val) => emit('update:visible', val)
})

const currentPath = ref('/')
const files = ref([])
const selectedFile = ref(null)
const fileInputRef = ref(null)
const transfers = ref([])
const renameDialogVisible = ref(false)
const mkdirDialogVisible = ref(false)
const newFileName = ref('')
const newDirName = ref('')
const renamingFile = ref(null)

// 计算属性
const pathParts = computed(() => {
    const parts = currentPath.value.split('/').filter(Boolean)
    let acc = ''
    return parts.map((part, index) => {
        acc += '/' + part
        return acc
    })
})

// 是否有已完成的传输
const hasCompletedTransfers = computed(() => {
    return transfers.value.some(t => t.status === 'completed')
})

// 监听 visible 变化
watch(() => props.visible, (val) => {
    if (val) {
        loadFiles()
    }
})

// 加载文件列表
const loadFiles = async () => {
    try {
        const data = await listFiles(props.sessionId, currentPath.value)
        if (data && data.data && data.data.files) {
            files.value = data.data.files
        } else {
            files.value = []
        }
    } catch (error) {
        console.error('加载文件列表失败:', error)
        files.value = []
        ElMessage.error('加载文件列表失败: ' + (error.message || error.response?.data?.message || error))
    }
}

// 行点击
const handleRowClick = (row) => {
    selectedFile.value = row
}

// 行双击
const handleRowDblClick = (row) => {
    if (row.isDir) {
        currentPath.value = row.path
        loadFiles()
    }
}

// 导航到指定路径
const navigateTo = (path) => {
    currentPath.value = path
    loadFiles()
}

// 刷新
const handleRefresh = () => {
    loadFiles()
}

// 上传点击
const handleUploadClick = () => {
    fileInputRef.value.click()
}

// 文件选择
const handleFileSelect = (event) => {
    const selectedFiles = event.target.files
    if (!selectedFiles || selectedFiles.length === 0) return

    Array.from(selectedFiles).forEach(file => {
        handleUploadFile(file)
    })

    event.target.value = '' // 重置 input
}

// 上传文件
const handleUploadFile = async (file) => {
    const transfer = {
        id: Date.now() + Math.random(),
        name: file.name,
        size: file.size,
        progress: 0,
        status: 'uploading',
        type: 'upload'
    }
    transfers.value.push(transfer)

    try {
        await uploadFile(
            props.sessionId,
            currentPath.value,
            file,
            (progress) => {
                transfer.progress = progress
            }
        )
        transfer.status = 'completed'
        ElMessage.success(`${file.name} 上传成功`)
        loadFiles()
    } catch (error) {
        transfer.status = 'failed'
        ElMessage.error(`${file.name} 上传失败: ` + error.message)
    }
}

// 下载点击
const handleDownloadClick = async () => {
    if (!selectedFile.value || selectedFile.value.isDir) {
        ElMessage.warning('请选择要下载的文件')
        return
    }

    try {
        await downloadFile(props.sessionId, selectedFile.value.path)
        ElMessage.success('下载已开始')
    } catch (error) {
        ElMessage.error('下载失败: ' + error.message)
    }
}

// 删除点击
const handleDeleteClick = async () => {
    if (!selectedFile.value) {
        ElMessage.warning('请选择要删除的文件')
        return
    }

    try {
        await ElMessageBox.confirm(
            `确定要删除 ${selectedFile.value.name} 吗？`,
            '确认删除',
            { type: 'warning' }
        )

        await deleteFile(props.sessionId, selectedFile.value.path)
        ElMessage.success('删除成功')
        loadFiles()
        selectedFile.value = null
    } catch (error) {
        if (error !== 'cancel') {
            ElMessage.error('删除失败: ' + error.message)
        }
    }
}

// 重命名点击
const handleRenameClick = (row) => {
    renamingFile.value = row
    newFileName.value = row.name
    renameDialogVisible.value = true
}

// 重命名确认
const handleRenameConfirm = async () => {
    if (!newFileName.value.trim()) {
        ElMessage.warning('请输入文件名')
        return
    }

    try {
        const newPath = renamingFile.value.path.substring(0, renamingFile.value.path.lastIndexOf('/') + 1) + newFileName.value
        await renameFile(props.sessionId, renamingFile.value.path, newPath)
        ElMessage.success('重命名成功')
        renameDialogVisible.value = false
        loadFiles()
    } catch (error) {
        ElMessage.error('重命名失败: ' + error.message)
    }
}

// 新建目录点击
const handleMkdirClick = () => {
    newDirName.value = ''
    mkdirDialogVisible.value = true
}

// 新建目录确认
const handleMkdirConfirm = async () => {
    if (!newDirName.value.trim()) {
        ElMessage.warning('请输入目录名')
        return
    }

    try {
        const newPath = currentPath.value.endsWith('/')
            ? currentPath.value + newDirName.value
            : currentPath.value + '/' + newDirName.value
        await createDir(props.sessionId, newPath)
        ElMessage.success('创建目录成功')
        mkdirDialogVisible.value = false
        loadFiles()
    } catch (error) {
        ElMessage.error('创建目录失败: ' + error.message)
    }
}

// 取消传输
const handleCancelTransfer = (transferId) => {
    transfers.value = transfers.value.filter(t => t.id !== transferId)
}

// 清除已完成的传输
const clearCompletedTransfers = () => {
    transfers.value = transfers.value.filter(t => t.status !== 'completed')
}

// 获取进度状态
const getProgressStatus = (status) => {
    switch (status) {
        case 'completed':
            return 'success'
        case 'failed':
            return 'exception'
        default:
            return ''
    }
}

// 关闭
const handleClose = () => {
    selectedFile.value = null
    currentPath.value = '/'
}

// 格式化文件大小
const formatSize = (bytes) => {
    if (bytes === 0) return '0 B'
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return (bytes / Math.pow(k, i)).toFixed(2) + ' ' + sizes[i]
}

// 格式化时间
const formatTime = (time) => {
    const date = new Date(time)
    return date.toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit'
    })
}

onMounted(() => {
    loadFiles()
})
</script>

<style scoped lang="scss">
.file-transfer-drawer {
    :deep(.el-drawer__header) {
        margin-bottom: 0px !important;
        padding-bottom: 8px !important;
    }

    :deep(.el-drawer__body) {
        padding: 12px;
        background-color: #f5f7fa;
    }
}

.file-transfer-container {
    height: 100%;
    display: flex;
    flex-direction: column;
    gap: 12px;
}

// 顶部路径栏
.header-section {
    flex-shrink: 0;
}

.path-input-wrapper {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px 12px;
    background-color: #fff;
    border-radius: 6px;
    border: 1px solid #dcdfe6;
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);

    .path-icon {
        color: #909399;
        font-size: 18px;
        flex-shrink: 0;
    }

    .path-breadcrumb {
        flex: 1;
        min-width: 0;

        :deep(.el-breadcrumb__item) {
            .el-breadcrumb__inner {
                cursor: pointer;
                font-size: 13px;
                color: #606266;
                transition: color 0.2s;

                &:hover {
                    color: #409eff;
                }
            }

            &:last-child {
                .el-breadcrumb__inner {
                    color: #303133;
                    font-weight: 500;
                }
            }
        }
    }
}

// 操作工具栏
.action-bar {
    flex-shrink: 0;
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 8px 12px;
    background-color: #fff;
    border-radius: 6px;
    border: 1px solid #dcdfe6;
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);

    .btn-group-main,
    .btn-group-secondary {
        display: flex;
        align-items: center;
    }

    .el-divider--vertical {
        height: 20px;
        margin: 0;
    }
}

// 文件列表
.file-list-wrapper {
    flex: 1;
    min-height: 0;
    background-color: #fff;
    border-radius: 6px;
    border: 1px solid #dcdfe6;
    overflow: hidden;
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);

    :deep(.el-table) {
        font-size: 13px;
        border: none;

        th.el-table__cell {
            background-color: #fafafa;
            color: #606266;
            font-weight: 600;
            padding: 8px 0;
            border-bottom: 1px solid #ebeef5;
        }

        td.el-table__cell {
            padding: 6px 0;
            border-bottom: 1px solid #ebeef5;
        }

        .el-table__row {
            cursor: pointer;
            transition: background-color 0.2s;

            &:hover {
                background-color: #f5f7fa !important;
            }

            &.current-row {
                background-color: #ecf5ff !important;
            }
        }

        .file-size {
            color: #909399;
            font-family: 'Consolas', monospace;
            font-size: 12px;
        }

        .parent-dir {
            color: #909399;
            font-weight: 500;
            font-style: italic;
        }
    }
}

// 传输进度区域
.transfer-section {
    flex-shrink: 0;
    max-height: 200px;
    background-color: #fff;
    border-radius: 6px;
    border: 1px solid #dcdfe6;
    overflow: hidden;
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
    display: flex;
    flex-direction: column;

    .transfer-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 8px 12px;
        background-color: #fafafa;
        border-bottom: 1px solid #ebeef5;

        .transfer-title {
            display: flex;
            align-items: center;
            gap: 6px;
            font-size: 13px;
            font-weight: 600;
            color: #303133;

            .el-icon {
                color: #409eff;
            }
        }
    }

    .transfer-list {
        flex: 1;
        overflow-y: auto;
        padding: 8px;
        max-height: 160px;

        &::-webkit-scrollbar {
            width: 4px;
        }

        &::-webkit-scrollbar-thumb {
            background: #dcdfe6;
            border-radius: 2px;
        }

        .transfer-item {
            padding: 8px;
            background-color: #f5f7fa;
            border-radius: 4px;
            margin-bottom: 6px;
            transition: all 0.2s;

            &:hover {
                background-color: #ecf5ff;
            }

            &:last-child {
                margin-bottom: 0;
            }
        }

        .transfer-info {
            display: flex;
            align-items: center;
            gap: 8px;
            margin-bottom: 8px;

            .transfer-icon {
                color: #409eff;
                font-size: 16px;
                flex-shrink: 0;
            }

            .transfer-name {
                flex: 1;
                font-size: 12px;
                color: #303133;
                overflow: hidden;
                text-overflow: ellipsis;
                white-space: nowrap;
            }
        }

        .transfer-status {
            display: flex;
            align-items: center;
            gap: 12px;

            .transfer-progress {
                flex: 1;
                min-width: 0;
            }

            .transfer-percent {
                font-size: 12px;
                color: #606266;
                font-family: 'Consolas', monospace;
                min-width: 40px;
                text-align: right;
            }
        }
    }
}
</style>

<style>
.file-transfer-drawer .el-drawer__header {
    margin-bottom: 0px !important;
    padding-bottom: 0px !important;
}
</style>
