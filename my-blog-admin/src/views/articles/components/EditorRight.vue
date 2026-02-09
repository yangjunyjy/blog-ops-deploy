<template>
    <el-card shadow="never" class="editor-right">
        <!-- 设置区 -->
        <div class="setting-section">
            <h3><el-icon>
                    <Setting />
                </el-icon> 文章设置</h3>

            <div class="setting-item">
                <label class="setting-label">封面图片</label>
                <div class="cover-uploader">
                    <div class="cover-preview" @click="triggerUpload" :style="{ backgroundImage: coverBg }">
                        <div v-if="!article.cover" class="upload-placeholder">
                            <el-icon>
                                <Plus />
                            </el-icon>
                            <span>点击上传封面</span>
                        </div>
                    </div>
                    <input ref="fileInput" type="file" accept="image/*" style="display: none"
                        @change="handleFileChange" />
                    <div class="upload-tips">
                        建议尺寸：1200×630px，支持 JPG、PNG 格式
                    </div>
                </div>
            </div>

            <div class="setting-item">
                <label class="setting-label">文章摘要</label>
                <el-input v-model="article.summary" type="textarea" :rows="4" maxlength="200" show-word-limit
                    placeholder="请输入文章摘要" />
            </div>

            <div class="setting-item">
                <label class="setting-label">URL 别名 (Slug)</label>
                <el-input v-model="article.slug" placeholder="例如: my-awesome-article"
                    :disabled="!!article.id">
                    <template #append>
                        <el-button @click="generateSlug">自动生成</el-button>
                    </template>
                </el-input>
                <div class="setting-tip">用于生成友好的 URL，留空则自动生成</div>
            </div>

            <div class="setting-item">
                <label class="setting-label">SEO 关键词</label>
                <el-input v-model="article.keywords" placeholder="多个关键词用逗号分隔" @keydown.enter="addKeyword" />
                <div class="keywords-list" v-if="keywordList.length">
                    <el-tag v-for="(keyword, index) in keywordList" :key="index" size="small" closable
                        @close="removeKeyword(index)">
                        {{ keyword }}
                    </el-tag>
                </div>
            </div>

            <div class="setting-item">
                <label class="setting-label">发布选项</label>
                <div class="publish-options">
                    <el-checkbox v-model="article.is_top" label="置顶文章" />
                    <el-checkbox v-model="article.is_recommended" label="推荐文章" />
                    <el-checkbox v-model="article.allow_comment" label="允许评论" />
                </div>
            </div>

            <div class="setting-item">
                <label class="setting-label">发布状态</label>
                <el-radio-group v-model="article.status">
                    <el-radio :value="0">草稿</el-radio>
                    <el-radio :value="1">发布</el-radio>
                    <el-radio :value="2">私密</el-radio>
                </el-radio-group>
            </div>

            <div class="setting-item">
                <label class="setting-label">发布时间</label>
                <el-date-picker v-model="article.publishedAt" type="datetime" placeholder="选择发布时间"
                    style="width: 100%" />
            </div>
        </div>

        <!-- 统计信息 -->
        <div class="stats-section">
            <h3><el-icon>
                    <DataLine />
                </el-icon> 文章统计</h3>
            <div class="stats-grid">
                <div class="stat-item">
                    <span class="stat-label">字数</span>
                    <span class="stat-value">{{ wordCount }}</span>
                </div>
                <div class="stat-item">
                    <span class="stat-label">图片数</span>
                    <span class="stat-value">{{ imageCount }}</span>
                </div>
                <div class="stat-item">
                    <span class="stat-label">阅读时间</span>
                    <span class="stat-value">{{ readingTime }}</span>
                </div>
                <div class="stat-item">
                    <span class="stat-label">段落数</span>
                    <span class="stat-value">{{ paragraphCount }}</span>
                </div>
            </div>
        </div>
    </el-card>
</template>

<script setup>
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { Setting, Plus, DataLine } from '@element-plus/icons-vue'
import { uploadImage as apiUploadImage } from '@/api/upload'

const props = defineProps({
    article: {
        type: Object,
        required: true
    },
    wordCount: {
        type: Number,
        default: 0
    },
    imageCount: {
        type: Number,
        default: 0
    },
    readingTime: {
        type: String,
        default: ''
    }
})

const emit = defineEmits(['cover-uploaded', 'update:article'])

const fileInput = ref(null)

// 计算属性
const coverBg = computed(() => {
    return props.article.cover ? `url(${props.article.cover})` : 'none'
})

const keywordList = computed({
    get() {
        return props.article.keywords ? props.article.keywords.split(',').filter(k => k.trim()) : []
    },
    set(val) {
        emit('update:article', { ...props.article, keywords: val.join(',') })
    }
})

const paragraphCount = computed(() => {
    return props.article.content.split('\n\n').filter(p => p.trim()).length
})

// 方法
const triggerUpload = () => {
    fileInput.value.click()
}

const handleFileChange = async (event) => {
    const file = event.target.files[0]
    if (!file) return

    // 验证文件类型
    if (!file.type.startsWith('image/')) {
        ElMessage.warning('请上传图片文件')
        return
    }

    // 验证文件大小（5MB）
    if (file.size > 5 * 1024 * 1024) {
        ElMessage.warning('图片大小不能超过5MB')
        return
    }

    try {
        const formData = new FormData()
        formData.append('file', file)

        const res = await apiUploadImage(formData)
        emit('cover-uploaded', res.data.url)
        ElMessage.success('封面图片上传成功')
    } catch (error) {
        ElMessage.error('图片上传失败')
    }

    // 清空input
    event.target.value = ''
}

const addKeyword = (event) => {
    const keyword = event.target.value.trim()
    if (keyword && !keywordList.value.includes(keyword)) {
        keywordList.value = [...keywordList.value, keyword]
        event.target.value = ''
    }
}

const removeKeyword = (index) => {
    keywordList.value.splice(index, 1)
}

// 生成 slug
const generateSlug = () => {
    if (!props.article.title) {
        ElMessage.warning('请先输入文章标题')
        return
    }
    // 将标题转换为 slug
    const slug = props.article.title
        .toLowerCase()
        .trim()
        .replace(/[\s\W-]+/g, '-') // 将非字母数字字符替换为连字符
        .replace(/^-+|-+$/g, '') // 移除首尾的连字符
        .substring(0, 100) // 限制长度

    emit('update:article', { ...props.article, slug })
    ElMessage.success('URL 别名已生成')
}
</script>

<style scoped lang="scss">
.editor-right {
    height: calc(100vh - 160px);
    overflow-y: auto;

    .setting-section {
        margin-bottom: 20px;

        h3 {
            display: flex;
            align-items: center;
            gap: 8px;
            margin: 0 0 20px;
            font-size: 16px;
            font-weight: 600;
            color: #303133;
        }

        .setting-item {
            margin-bottom: 20px;

            .setting-label {
                display: block;
                margin-bottom: 8px;
                font-size: 14px;
                font-weight: 500;
                color: #606266;
            }

            .cover-uploader {
                .cover-preview {
                    width: 100%;
                    height: 160px;
                    background-color: #f5f7fa;
                    background-size: cover;
                    background-position: center;
                    border-radius: 4px;
                    border: 1px dashed #dcdfe6;
                    cursor: pointer;
                    display: flex;
                    align-items: center;
                    justify-content: center;
                    transition: border-color 0.3s;

                    &:hover {
                        border-color: #409eff;

                        .upload-placeholder {
                            color: #409eff;
                        }
                    }

                    .upload-placeholder {
                        display: flex;
                        flex-direction: column;
                        align-items: center;
                        color: #909399;
                        transition: color 0.3s;

                        .el-icon {
                            font-size: 24px;
                            margin-bottom: 8px;
                        }
                    }
                }

                .upload-tips {
                    font-size: 12px;
                    color: #909399;
                    margin-top: 8px;
                }
            }

            .keywords-list {
                margin-top: 10px;

                .el-tag {
                    margin-right: 8px;
                    margin-bottom: 8px;
                }
            }

            .setting-tip {
                font-size: 12px;
                color: #909399;
                margin-top: 6px;
            }

            .publish-options {
                .el-checkbox {
                    display: block;
                    margin-bottom: 8px;
                }
            }
        }
    }

    .stats-section {
        background: #f5f7fa;
        padding: 20px;
        border-radius: 8px;
        margin-bottom: 20px;

        h3 {
            display: flex;
            align-items: center;
            gap: 8px;
            margin: 0 0 15px;
            font-size: 16px;
            font-weight: 600;
        }

        .stats-grid {
            .stat-item {
                display: flex;
                justify-content: space-between;
                padding: 10px 0;
                border-bottom: 1px solid #e4e7ed;

                &:last-child {
                    border-bottom: none;
                }

                .stat-label {
                    color: #909399;
                    font-size: 14px;
                }

                .stat-value {
                    font-weight: 600;
                    color: #409eff;
                    font-size: 14px;
                }
            }
        }
    }

    .action-section {
        .el-button {
            margin-bottom: 10px;
        }
    }
}
</style>