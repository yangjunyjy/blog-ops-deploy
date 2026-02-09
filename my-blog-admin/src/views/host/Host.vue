<template>
    <div class="host-list">
        <el-card shadow="never">
            <!-- 搜索栏 -->
            <el-form :model="searchForm" inline class="search-form" @submit.prevent>
                <div class="form-items">
                    <el-form-item label="主机名称">
                        <el-input v-model="searchForm.name" placeholder="主机名称" clearable style="width: 150px" />
                    </el-form-item>
                    <el-form-item label="主机IP">
                        <el-input v-model="searchForm.ip" placeholder="主机IP" clearable style="width: 150px" />
                    </el-form-item>
                    <el-form-item label="认证类型">
                        <el-select v-model="searchForm.auth_type" placeholder="认证类型" style="width: 120px" clearable
                            :teleported="false">
                            <el-option label="密码" value="password" />
                            <el-option label="密钥" value="key" />
                        </el-select>
                    </el-form-item>
                    <el-select v-model="searchForm.status" placeholder="状态" style="width: 80px" clearable
                        :teleported="false">
                        <el-option label="启用" :value="1" />
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
                    <span class="action-title">主机管理</span>
                    <div class="action-buttons">
                        <HaloButton type="primary" size="medium" :icon="Plus" content="新增" @click="handleCreate" />
                    </div>
                </div>

                <el-table :data="hostList" v-loading="loading" :header-cell-style="{ backgroundColor: '#f5f7fa' }"
                    row-key="id">
                    <el-table-column prop="id" label="ID" width="70" align="center" />
                    <el-table-column prop="name" label="主机名称" min-width="120" />
                    <el-table-column prop="ip" label="主机IP" min-width="120" />
                    <el-table-column prop="port" label="端口" width="80" align="center" />
                    <el-table-column prop="auth_type" label="认证类型" width="100" align="center">
                        <template #default="{ row }">
                            <el-tag :type="row.auth_type === 'password' ? 'primary' : 'success'" size="small">
                                {{ row.auth_type === 'password' ? '密码' : '密钥' }}
                            </el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column label="状态" width="80">
                        <template #default="{ row }">
                            <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
                                {{ row.status === 1 ? '启用' : '禁用' }}
                            </el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column prop="description" label="描述" min-width="200" />
                    <el-table-column label="操作" width="240" fixed="right">
                        <template #default="{ row }">
                            <div style="display: flex; align-items: center; gap: 8px;">
                                <el-button type="primary" link size="small" @click="handleRemoteConnect(row)"
                                    class="icon-btn" :title="`使用 ${getConnectTypeName(row.connect_type)} 连接`">
                                    <el-icon class="action-icon">
                                        <img src="/svg/connect.svg" alt="connect" style="width: 16px; height: 16px;" />
                                    </el-icon>
                                    <span style="margin-left: 4px;">{{ getConnectTypeName(row.connect_type) }}</span>
                                </el-button>
                                <!-- 文件传输功能暂未实现 -->
                                <!-- <el-button type="success" link size="small" @click="handleTransferFile(row)"
                                    class="icon-btn" title="文件传输">
                                    <el-icon class="action-icon">
                                        <img src="/svg/transfer.svg" alt="transfer"
                                            style="width: 16px; height: 16px;" />
                                    </el-icon>
                                    <span style="margin-left: 4px;">传输</span>
                                </el-button> -->
                                <el-button type="warning" link size="small" @click="handleEdit(row)" class="icon-btn">
                                    <el-icon class="action-icon">
                                        <img src="/svg/edit.svg" alt="edit" style="width: 16px; height: 16px;" />
                                    </el-icon>
                                </el-button>
                                <el-button type="danger" link size="small" @click="handleDelete(row.id)"
                                    class="icon-btn">
                                    <el-icon class="action-icon delete-icon">
                                        <img src="/svg/delete.svg" alt="delete"
                                            style="width: 16px; height: 16px; color: #f56c6c;" class="svg-red" />
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
                    :page-sizes="[10, 20, 50, 100]" :total="pagination.total"
                    layout="total, sizes, prev, pager, next, jumper" @size-change="handleSizeChange"
                    @current-change="handleCurrentChange" />
            </div>
        </el-card>

        <!-- 编辑对话框 -->
        <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px">
            <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
                <el-form-item label="主机名称" prop="name">
                    <el-input v-model="form.name" placeholder="请输入主机名称" />
                </el-form-item>
                <el-form-item label="主机IP" prop="ip">
                    <el-input v-model="form.ip" placeholder="请输入主机IP地址" />
                </el-form-item>
                <el-form-item label="端口" prop="port">
                    <el-input v-model="form.port" placeholder="默认22" />
                </el-form-item>
                <el-form-item label="认证类型" prop="auth_type">
                    <el-radio-group v-model="form.auth_type" @change="handleAuthTypeChange">
                        <el-radio label="password">密码</el-radio>
                        <el-radio label="key">密钥</el-radio>
                    </el-radio-group>
                </el-form-item>
                <el-form-item v-if="form.auth_type === 'password'" label="用户名" prop="username">
                    <el-input v-model="form.username" placeholder="请输入用户名" />
                </el-form-item>
                <el-form-item v-if="form.auth_type === 'password'" label="密码" prop="password">
                    <el-input v-model="form.password" type="password" placeholder="请输入密码" show-password />
                </el-form-item>
                <el-form-item v-if="form.auth_type === 'key'" label="私钥" prop="private_key">
                    <el-input v-model="form.private_key" type="textarea" :rows="4" placeholder="请输入SSH私钥" />
                </el-form-item>
                <el-form-item label="连接方式" prop="connect_type">
                    <el-select v-model="form.connect_type" placeholder="请选择连接方式" style="width: 100%">
                        <el-option label="Web 终端" value="web" />
                        <el-option label="Xshell" value="xshell" />
                        <el-option label="SecureCRT" value="securecrt" />
                        <el-option label="PuTTY" value="putty" />
                    </el-select>
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

        <!-- 文件传输对话框 -->
        <!-- <FileTransferClient v-model:visible="fileTransferDialogVisible" :host="currentHost"
            @close="handleFileTransferClose" /> -->
    </div>
</template>

<script setup>
defineOptions({
    name: 'Host'
})

import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import HaloButton from '../../layout/components/HaloButton.vue'
import {
    getHostList,
    createHost,
    updateHost,
    deleteHost
} from '@/api/host'

const loading = ref(false)
const dialogVisible = ref(false)
// const fileTransferDialogVisible = ref(false) // 文件传输功能暂未实现
const submitting = ref(false)
const formRef = ref(null)
const hostList = ref([])
// const currentHost = ref(null) // 文件传输功能暂未实现

// 分页
const pagination = reactive({
    page: 1,
    pageSize: 10,
    total: 0
})



const searchForm = reactive({
    name: '',
    ip: '',
    auth_type: '',
    status: 1
})

const dialogTitle = ref('新增主机')
const form = reactive({
    id: null,
    name: '',
    ip: '',
    port: 22,
    auth_type: 'password',
    username: '',
    password: '',
    private_key: '',
    connect_type: 'web',
    description: '',
    status: 1
})





const rules = {
    name: [{ required: true, message: '请输入主机名称', trigger: 'blur' }],
    ip: [{ required: true, message: '请输入主机IP', trigger: 'blur' }],
    port: [{ required: true, message: '请输入端口', trigger: 'blur' }],
    username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
    password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
    private_key: [{ required: true, message: '请输入私钥', trigger: 'blur' }]
}

const loadList = async () => {
    loading.value = true
    try {
        const params = {
            page: pagination.page,
            page_size: pagination.pageSize,
            name: searchForm.name || '',
            address: searchForm.ip || '',
            type: searchForm.auth_type || '',
            status: searchForm.status === 1 ? 'active' : searchForm.status === 0 ? 'inactive' : ''
        }
        const data = await getHostList(params)
        // 后端返回格式: {code: 200, message: '', data: {total, items}}
        if (data && data.data) {
            pagination.total = data.data.total || 0
            // 使用适配器转换数据格式
            hostList.value = (data.data.items || []).map(host => ({
                ...host,
                ip: host.address,          // address -> ip
                auth_type: host.type,       // type -> auth_type
                connect_type: 'web',         // 默认web连接
                description: ''
            }))
        }
    } catch (error) {
        console.error('加载主机列表失败:', error)
        ElMessage.error('加载失败')
    } finally {
        loading.value = false
    }
}

const handleSearch = () => {
    pagination.page = 1
    // 如果后端支持搜索，直接调用loadList
    loadList()
}

const handleReset = () => {
    searchForm.name = ''
    searchForm.ip = ''
    searchForm.auth_type = ''
    pagination.page = 1
    searchForm.status = 1
    loadList()
}

const handleCreate = () => {
    dialogTitle.value = '新增主机'
    Object.assign(form, {
        id: null,
        name: '',
        ip: '',
        port: 22,
        auth_type: 'password',
        username: '',
        password: '',
        private_key: '',
        connect_type: 'web',
        description: '',
        status: 1
    })
    dialogVisible.value = true
}

const handleEdit = (row) => {
    dialogTitle.value = '编辑主机'
    Object.assign(form, {
        id: row.id,
        name: row.name,
        ip: row.ip,
        port: row.port,
        auth_type: row.auth_type,
        username: row.username || '',
        password: row.password || '',
        private_key: row.private_key || '',
        connect_type: row.connect_type || 'web',
        description: row.description,
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
                // 更新主机
                const params = {
                    id: form.id,
                    name: form.name,
                    address: form.ip,              // ip -> address
                    port: Number(form.port),
                    username: form.username,
                    password: form.password,
                    secret_key: form.private_key,    // private_key -> secret_key
                    type: form.auth_type,           // auth_type -> type
                    status: form.status === 1 ? 'active' : 'inactive'  // 1/0 -> 'active'/'inactive'
                }
                await updateHost(params)
                // 响应拦截器已处理，直接显示成功
                ElMessage.success('更新成功')
            } else {
                // 创建主机
                const params = {
                    name: form.name,
                    address: form.ip,
                    port: Number(form.port),
                    username: form.username,
                    password: form.password,
                    secret_key: form.private_key,
                    type: form.auth_type,
                    status: form.status === 1 ? 'active' : 'inactive'
                }
                await createHost(params)
                ElMessage.success('创建成功')
            }
            dialogVisible.value = false
            loadList()
        } catch (error) {
            console.error('操作失败:', error)
            // 错误已由响应拦截器处理
        } finally {
            submitting.value = false
        }
    })
}

const handleDelete = async (id) => {
    try {
        await ElMessageBox.confirm('确定删除该主机吗？', '提示', { type: 'warning' })
        await deleteHost(id)
        ElMessage.success('删除成功')
        loadList()
    } catch (error) {
        if (error !== 'cancel') {
            ElMessage.error('删除失败')
        }
    }
}

const handleRemoteConnect = (row) => {
    const connectType = row.connect_type || 'web'

    console.log('点击连接，主机信息:', row)

    if (connectType === 'web') {
        // Web 终端：在新标签页打开多标签会话页面
        const url = window.location.origin + '/terminal/' + row.id
        console.log('打开终端 URL:', url, '主机 ID:', row.id)
        window.open(url, '_blank')
    } else if (connectType === 'xshell') {
        // Xshell：生成 xshell 连接配置并下载
        generateXshellConfig(row)
    } else if (connectType === 'securecrt') {
        // SecureCRT：生成 securecrt 连接配置并下载
        generateSecureCRTConfig(row)
    } else if (connectType === 'putty') {
        // PuTTY：生成 putty 注册表文件并下载
        generatePuttyConfig(row)
    }
}

// const handleTransferFile = (row) => {
//     currentHost.value = row
//     fileTransferDialogVisible.value = true
// }

// const handleFileTransferClose = () => {
//     // 清理操作，如果需要的话
// }

// 生成 Xshell 连接配置
const generateXshellConfig = (host) => {
    const authInfo = host.auth_type === 'password'
        ? `Password=${host.password}`
        : `PublicKeyFile=C:\\Users\\${encodeURIComponent(host.username)}\\.ssh\\id_rsa`

    const config = `[ConnectInfo]
Host=${host.ip}
Port=${host.port}
UserName=${host.username}
${authInfo}

[Terminal]
TerminalType=xterm

[Font]
FontName=Consolas
FontSize=10
`
    downloadConfig(host, config, 'xsh')
}

// 生成 SecureCRT 连接配置
const generateSecureCRTConfig = (host) => {
    const config = `S:"Hostname"="${host.ip}"
S:"Username"="${host.username}"
D:"Port [SSH]"=${host.port}
D:"Authentication Method"=2
${host.auth_type === 'password' ? `S:"Password"="${host.password}"` : 'D:"Use Public Key"=1'}
S:"Session Name"="${host.name}"
`
    downloadConfig(host, config, 'ini')
}

// 生成 PuTTY 连接配置
const generatePuttyConfig = (host) => {
    const authInfo = host.auth_type === 'password'
        ? `S:"Password"="${host.password}"`
        : ''

    const config = `Windows Registry Editor Version 5.00

[HKEY_CURRENT_USER\\Software\\SimonTatham\\PuTTY\\Sessions\\${host.name}]
"HostName"="${host.ip}"
"PortNumber"=dword:${host.port}
"UserName"="${host.username}"
"Protocol"="ssh"
${authInfo}
`
    downloadConfig(host, config, 'reg')
}

// 下载配置文件
const downloadConfig = (host, content, ext) => {
    const blob = new Blob([content], { type: 'text/plain;charset=utf-8' })
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `${host.name}.${ext}`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    URL.revokeObjectURL(url)
    ElMessage.success(`${host.name} 连接配置已下载，请导入对应客户端`)
}

const handleAuthTypeChange = (val) => {
    // 切换认证类型时清空另一个字段
    if (val === 'password') {
        form.private_key = ''
    } else {
        form.username = ''
        form.password = ''
    }
}

const handleSizeChange = (size) => {
    pagination.pageSize = size
    pagination.page = 1
    loadList()
}

const handleCurrentChange = (page) => {
    pagination.page = page
    loadList()
}

// 获取连接方式名称
const getConnectTypeName = (type) => {
    const typeMap = {
        web: 'Web终端',
        xshell: 'Xshell',
        securecrt: 'SecureCRT',
        putty: 'PuTTY'
    }
    return typeMap[type] || 'Web终端'
}

onMounted(() => {
    loadList()
})
</script>

<style scoped lang="scss">
.host-list {
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