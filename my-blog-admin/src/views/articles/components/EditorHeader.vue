<template>
    <el-card shadow="never" class="editor-header">
        <div class="header-content">
            <div class="header-left">
                <HaloButton content="返回" size="medium" type="default" @click="$emit('back')">
                    <template #icon>
                        <el-icon><ArrowLeft /></el-icon>
                    </template>
                </HaloButton>
                <div class="header-title">
                    <h2 v-if="articleId">编辑文章</h2>
                    <h2 v-else>写文章</h2>
                    <span class="article-id" v-if="articleId">#{{ articleId }}</span>
                </div>
                <el-upload ref="uploadRef" :show-file-list="false" accept=".md,.markdown" :before-upload="beforeUpload"
                    :http-request="handleFileUpload" :on-success="handleUploadSuccess" :on-error="handleUploadError"
                    :on-progress="handleUploadProgress" :disabled="uploading" class="import-upload">
                    <HaloButton content="导入 Markdown" size="medium" type="success" :loading="uploading">
                        <template #icon>
                            <el-icon v-if="!uploading"><Upload /></el-icon>
                        </template>
                    </HaloButton>
                </el-upload>
            </div>

            <div class="header-right">
                <HaloButton content="保存草稿" size="medium" type="default" @click="$emit('save-draft')">
                    <template #icon>
                        <el-icon><Document /></el-icon>
                    </template>
                </HaloButton>
                <HaloButton content="预览" size="medium" type="warning" @click="$emit('preview')">
                    <template #icon>
                        <el-icon><View /></el-icon>
                    </template>
                </HaloButton>
                <HaloButton :content="publishing ? '发布中...' : '发布文章'" size="medium" type="primary"
                    :loading="publishing" @click="$emit('publish')">
                    <template #icon>
                        <el-icon v-if="!publishing"><Promotion /></el-icon>
                    </template>
                </HaloButton>
            </div>
        </div>
    </el-card>
</template>

<script setup>
import { ref, defineProps, defineEmits } from 'vue'
import { ElMessage } from 'element-plus'
import { ArrowLeft, Document, View, Promotion, Upload } from '@element-plus/icons-vue'
import { uploadMarkdown } from '@/api/upload'
import HaloButton from '@/layout/components/HaloButton.vue'

defineProps({
    publishing: {
        type: Boolean,
        default: false
    },
    articleId: {
        type: [String, Number],
        default: null
    }
})

const emit = defineEmits(['back', 'preview', 'publish', 'import'])

const uploading = ref(false)
const uploadRef = ref()

// 文件上传前的校验
const beforeUpload = (file) => {
    const isMarkdown = file.name.endsWith('.md') || file.name.endsWith('.markdown')
    const isLt5M = file.size / 1024 / 1024 < 5

    if (!isMarkdown) {
        ElMessage.error('只支持 .md 或 .markdown 格式的文件！')
        return false
    }
    if (!isLt5M) {
        ElMessage.error('文件大小不能超过 5MB！')
        return false
    }
    return true
}

// 自定义上传
const handleFileUpload = async (options) => {
    const { file } = options
    uploading.value = true

    // 创建 FormData
    const formData = new FormData()
    formData.append('file', file)

    try {
        const res = await uploadMarkdown(formData)

        // 根据后端返回的数据结构进行调整
        if (res.code === 200) {
            // 假设后端返回的数据结构为：
            // {
            //   code: 200,
            //   data: {
            //     title: '解析出的标题',
            //     content: '文件内容',
            //     summary: '摘要',
            //     // 其他可能的信息
            //   }
            // }
            const data = res.data || {}

            // 触发导入事件，将解析的数据传递给父组件
            emit('import', {
                title: data.title || file.name.replace(/\.(md|markdown)$/, ''),
                content: data.content || '',
                summary: data.summary || '',
                // 可以传递其他解析出的信息，如分类、标签等
                ...data
            })

            ElMessage.success('文件导入成功')
        } else {
            ElMessage.error(res.msg || '文件解析失败')
        }
    } catch (error) {
        console.error('上传失败:', error)
        ElMessage.error('上传失败，请稍后重试')
    } finally {
        uploading.value = false
    }
}

// 上传成功回调
const handleUploadSuccess = (response, file) => {
    // 这里我们已经在 handleFileUpload 中处理了，所以可以留空
    // 或者也可以在这里处理，根据你的需求
}

// 上传失败回调
const handleUploadError = (error, file) => {
    uploading.value = false
    console.error('上传失败:', error)
    ElMessage.error('上传失败，请检查网络连接')
}

// 上传进度回调
const handleUploadProgress = (event, file) => {
    // 如果需要显示上传进度，可以在这里处理
    // console.log('上传进度:', event.percent)
}
</script>

<style scoped lang="scss">
.editor-header {
    margin-bottom: 0;
    border-radius: 0;
    border-left: none;
    border-right: none;
    border-top: none;
    position: sticky;
    top: 0;
    z-index: 100;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);

    .header-content {
        display: flex;
        justify-content: space-between;
        align-items: center;

        .header-left {
            display: flex;
            align-items: center;
            gap: 20px;

            .header-title {
                display: flex;
                align-items: center;
                gap: 10px;

                h2 {
                    margin: 0;
                    font-size: 18px;
                    color: #303133;
                    font-weight: 600;
                }

                .article-id {
                    font-size: 14px;
                    color: #909399;
                    background: #f5f7fa;
                    padding: 2px 8px;
                    border-radius: 2px;
                }
            }

            .import-upload {
                :deep(.el-upload) {
                    display: inline-block;
                }
            }
        }

        .header-right {
            display: flex;
            align-items: center;
            gap: 10px;
        }
    }
}
</style>
