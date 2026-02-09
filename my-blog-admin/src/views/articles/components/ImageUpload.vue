<template>
    <div class="image-upload">
        <div class="upload-methods">
            <!-- 上传文件 -->
            <div class="upload-method">
                <div class="upload-area" @click="triggerFileUpload">
                    <el-icon>
                        <Upload />
                    </el-icon>
                    <span>本地上传</span>
                    <input ref="fileInput" type="file" accept="image/*" multiple style="display: none"
                        @change="handleFileUpload" />
                </div>
            </div>

            <!-- 粘贴图片 -->
            <div class="upload-method">
                <div class="upload-area" @click="handlePasteClick">
                    <el-icon>
                        <DocumentCopy />
                    </el-icon>
                    <span>粘贴图片</span>
                    <input ref="pasteInput" type="text" placeholder="Ctrl+V 粘贴图片" @paste="handlePaste"
                        style="display: none" />
                </div>
            </div>

            <!-- 网络图片 -->
            <div class="upload-method">
                <div class="upload-area" @click="showUrlInput = true">
                    <el-icon>
                        <Link />
                    </el-icon>
                    <span>网络图片</span>
                </div>
            </div>
        </div>

        <!-- 图片列表 -->
        <div class="image-list" v-if="images.length">
            <div class="image-item" v-for="(image, index) in images" :key="index">
                <img :src="image.url" :alt="image.name" @load="image.loaded = true" />
                <div class="image-info">
                    <span class="image-name">{{ image.name }}</span>
                    <span class="image-size">{{ formatSize(image.size) }}</span>
                </div>
                <div class="image-actions">
                    <el-button size="small" @click="selectImage(image)">选择</el-button>
                    <el-button size="small" type="danger" @click="removeImage(index)">
                        删除
                    </el-button>
                </div>
            </div>
        </div>

        <!-- 网络图片输入 -->
        <el-dialog v-model="showUrlInput" title="添加网络图片" width="400px">
            <div class="url-input">
                <el-input v-model="imageUrl" placeholder="请输入图片URL" @keydown.enter="addUrlImage" />
                <div class="url-tips">
                    支持 JPG、PNG、GIF、WebP 格式
                </div>
                <div class="dialog-actions">
                    <el-button @click="showUrlInput = false">取消</el-button>
                    <el-button type="primary" @click="addUrlImage">确认</el-button>
                </div>
            </div>
        </el-dialog>
    </div>
</template>

<script setup>
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { Upload, DocumentCopy, Link } from '@element-plus/icons-vue'
import { uploadImage as apiUploadImage } from '@/api/upload'

const emit = defineEmits(['upload-success', 'cancel'])

const fileInput = ref(null)
const pasteInput = ref(null)
const showUrlInput = ref(false)
const images = ref([])
const imageUrl = ref('')

// 触发文件选择
const triggerFileUpload = () => {
    fileInput.value.click()
}

// 处理文件上传
const handleFileUpload = async (event) => {
    const files = Array.from(event.target.files)

    for (const file of files) {
        if (!file.type.startsWith('image/')) {
            ElMessage.warning(`${file.name} 不是图片文件`)
            continue
        }

        // 验证文件大小（10MB）
        if (file.size > 10 * 1024 * 1024) {
            ElMessage.warning(`${file.name} 大小不能超过10MB`)
            continue
        }

        try {
            const formData = new FormData()
            formData.append('file', file)

            const res = await apiUploadImage(formData)

            images.value.push({
                url: res.data.url,
                name: file.name,
                size: file.size,
                loaded: false
            })

            ElMessage.success(`${file.name} 上传成功`)
        } catch (error) {
            ElMessage.error(`${file.name} 上传失败`)
        }
    }

    // 清空input
    event.target.value = ''
}

// 处理粘贴
const handlePasteClick = () => {
    pasteInput.value.focus()
    ElMessage.info('请按 Ctrl+V 粘贴图片')
}

const handlePaste = async (event) => {
    const items = event.clipboardData?.items
    if (!items) return

    for (const item of items) {
        if (item.type.indexOf('image') !== -1) {
            event.preventDefault()
            const file = item.getAsFile()
            if (file) {
                await uploadPasteImage(file)
            }
        }
    }
}

const uploadPasteImage = async (file) => {
    try {
        const formData = new FormData()
        formData.append('file', file)

        const res = await apiUploadImage(formData)

        images.value.push({
            url: res.data.url,
            name: '粘贴图片.png',
            size: file.size,
            loaded: false
        })

        ElMessage.success('图片粘贴成功')
    } catch (error) {
        ElMessage.error('图片粘贴失败')
    }
}

// 添加网络图片
const addUrlImage = () => {
    if (!imageUrl.value.trim()) {
        ElMessage.warning('请输入图片URL')
        return
    }

    // 验证URL格式
    const urlPattern = /^(https?:\/\/).+\.(jpg|jpeg|png|gif|webp)(\?.*)?$/i
    if (!urlPattern.test(imageUrl.value)) {
        ElMessage.warning('请输入有效的图片URL')
        return
    }

    images.value.push({
        url: imageUrl.value,
        name: '网络图片',
        size: 0,
        loaded: false
    })

    showUrlInput.value = false
    imageUrl.value = ''
    ElMessage.success('网络图片添加成功')
}

// 选择图片
const selectImage = (image) => {
    emit('upload-success', image.url)
}

// 删除图片
const removeImage = (index) => {
    images.value.splice(index, 1)
}

// 格式化文件大小
const formatSize = (bytes) => {
    if (bytes === 0) return '未知大小'
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}
</script>

<style scoped lang="scss">
.image-upload {
    .upload-methods {
        display: grid;
        grid-template-columns: repeat(3, 1fr);
        gap: 10px;
        margin-bottom: 20px;

        .upload-method {
            .upload-area {
                height: 100px;
                border: 2px dashed #dcdfe6;
                border-radius: 8px;
                display: flex;
                flex-direction: column;
                align-items: center;
                justify-content: center;
                cursor: pointer;
                transition: all 0.3s;

                &:hover {
                    border-color: #409eff;
                    background: #ecf5ff;

                    .el-icon,
                    span {
                        color: #409eff;
                    }
                }

                .el-icon {
                    font-size: 28px;
                    color: #909399;
                    margin-bottom: 10px;
                }

                span {
                    color: #606266;
                    font-size: 14px;
                }
            }
        }
    }

    .image-list {
        .image-item {
            display: flex;
            align-items: center;
            padding: 10px;
            border: 1px solid #e4e7ed;
            border-radius: 4px;
            margin-bottom: 10px;

            img {
                width: 60px;
                height: 60px;
                object-fit: cover;
                border-radius: 4px;
                margin-right: 15px;
            }

            .image-info {
                flex: 1;
                display: flex;
                flex-direction: column;

                .image-name {
                    font-size: 14px;
                    color: #303133;
                    margin-bottom: 5px;
                }

                .image-size {
                    font-size: 12px;
                    color: #909399;
                }
            }

            .image-actions {
                .el-button {
                    margin-left: 5px;
                }
            }
        }
    }

    .url-input {
        .url-tips {
            font-size: 12px;
            color: #909399;
            margin-top: 8px;
        }

        .dialog-actions {
            margin-top: 20px;
            text-align: right;
        }
    }
}
</style>