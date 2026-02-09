<template>
    <el-dialog v-model="dialogVisible" title="文章预览" fullscreen class="preview-dialog">
        <template #header>
            <div class="preview-header">
                <span class="preview-title">文章预览</span>
                <div class="preview-actions">
                    <el-button @click="$emit('update:modelValue', false)">
                        关闭
                    </el-button>
                    <el-button type="primary" @click="handlePrint">
                        打印
                    </el-button>
                    <el-button type="success" @click="handleExport">
                        导出PDF
                    </el-button>
                </div>
            </div>
        </template>

        <div class="preview-content" ref="previewRef">
            <!-- 文章头部 -->
            <div class="article-header">
                <h1 class="article-title">{{ article.title }}</h1>
                <div class="article-meta">
                    <span class="meta-item">
                        <el-icon>
                            <Calendar />
                        </el-icon>
                        {{ formatDate(article.createdAt) }}
                    </span>
                    <span class="meta-item">
                        <el-icon>
                            <Folder />
                        </el-icon>
                        {{ article.category?.name || '未分类' }}
                    </span>
                    <span class="meta-item">
                        <el-icon>
                            <Clock />
                        </el-icon>
                        {{ readingTime }}
                    </span>
                    <span class="meta-item">
                        <el-icon>
                            <View />
                        </el-icon>
                        阅读量: 0
                    </span>
                </div>
            </div>

            <!-- 文章封面 -->
            <div class="article-cover" v-if="article.cover">
                <img :src="article.cover" alt="文章封面" />
            </div>

            <!-- 文章内容 -->
            <div class="article-body">
                <div class="markdown-preview" v-html="renderedContent"></div>
            </div>

            <!-- 文章底部 -->
            <div class="article-footer">
                <div class="article-tags">
                    <el-tag v-for="tag in article.tags" :key="tag.id" size="medium" type="info">
                        {{ tag.name }}
                    </el-tag>
                </div>

                <div class="article-actions">
                    <div class="action-buttons">
                        <el-button type="primary" icon="Star" circle />
                        <el-button type="success" icon="Share" circle />
                        <el-button type="warning" icon="ChatDotRound" circle />
                    </div>
                </div>
            </div>
        </div>
    </el-dialog>
</template>

<script setup>
import { ref, computed, watch, nextTick } from 'vue'
import { Calendar, Folder, Clock, View, Star, Share, ChatDotRound } from '@element-plus/icons-vue'

const props = defineProps({
    modelValue: Boolean,
    article: Object,
    renderedContent: String
})

const emit = defineEmits(['update:modelValue'])

const dialogVisible = computed({
    get() {
        return props.modelValue
    },
    set(val) {
        emit('update:modelValue', val)
    }
})

const previewRef = ref(null)

// 计算阅读时间
const readingTime = computed(() => {
    const wordCount = props.article.content?.replace(/[^\u4e00-\u9fa5\w]/g, '').length || 0
    const minutes = Math.ceil(wordCount / 300)
    return `${minutes} 分钟阅读`
})

// 格式化日期
const formatDate = (date) => {
    if (!date) return ''
    const d = new Date(date)
    return d.toLocaleString('zh-CN', {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    })
}

// 打印文章
const handlePrint = () => {
    window.print()
}

// 导出PDF（模拟）
const handleExport = () => {
    // 这里可以集成PDF导出库，如html2pdf.js
    alert('PDF导出功能需要集成html2pdf.js')
}

// 监听内容变化，添加代码高亮
watch(() => props.renderedContent, () => {
    nextTick(() => {
        highlightCodeBlocks()
    })
}, { immediate: true })

const highlightCodeBlocks = () => {
    if (!previewRef.value) return

    const codeBlocks = previewRef.value.querySelectorAll('pre code')
    codeBlocks.forEach(block => {
        // 这里可以集成highlight.js
        block.classList.add('hljs')
    })
}
</script>

<style scoped lang="scss">
.preview-dialog {
    :deep(.el-dialog__body) {
        padding: 0;
    }

    .preview-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        width: 100%;

        .preview-title {
            font-size: 18px;
            font-weight: 600;
        }

        .preview-actions {
            .el-button {
                margin-left: 10px;
            }
        }
    }

    .preview-content {
        max-width: 800px;
        margin: 0 auto;
        padding: 40px 20px;

        .article-header {
            text-align: center;
            margin-bottom: 40px;

            .article-title {
                font-size: 32px;
                font-weight: 700;
                margin: 0 0 20px;
                color: #1a1a1a;
                line-height: 1.3;
            }

            .article-meta {
                display: flex;
                justify-content: center;
                flex-wrap: wrap;
                gap: 20px;
                color: #666;
                font-size: 14px;

                .meta-item {
                    display: flex;
                    align-items: center;
                    gap: 5px;

                    .el-icon {
                        font-size: 16px;
                    }
                }
            }
        }

        .article-cover {
            margin-bottom: 40px;

            img {
                width: 100%;
                max-height: 400px;
                object-fit: cover;
                border-radius: 8px;
            }
        }

        .article-body {
            .markdown-preview {
                font-size: 16px;
                line-height: 1.8;
                color: #333;

                h1,
                h2,
                h3,
                h4,
                h5,
                h6 {
                    margin: 30px 0 15px;
                    color: #1a1a1a;
                    font-weight: 600;
                }

                h1 {
                    font-size: 28px;
                }

                h2 {
                    font-size: 24px;
                }

                h3 {
                    font-size: 20px;
                }

                h4 {
                    font-size: 18px;
                }

                p {
                    margin-bottom: 20px;
                }

                a {
                    color: #0366d6;
                    text-decoration: none;

                    &:hover {
                        text-decoration: underline;
                    }
                }

                img {
                    max-width: 100%;
                    height: auto;
                    border-radius: 4px;
                }

                pre {
                    background: #282c34;
                    border-radius: 8px;
                    padding: 0;
                    overflow-x: auto;
                    margin: 20px 0;
                    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);

                    code {
                        font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
                        font-size: 14px;
                        line-height: 1.6;
                        background: transparent;
                        padding: 16px;
                        display: block;
                    }
                }

                .code-block-wrapper {
                    background-color: #282c34;
                    border-radius: 8px;
                    margin: 16px 0;
                    overflow: hidden;
                    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);

                    .code-header {
                        display: flex;
                        justify-content: space-between;
                        align-items: center;
                        padding: 8px 16px;
                        background-color: #21252b;
                        border-bottom: 1px solid #3e4451;

                        .code-language {
                            color: #abb2bf;
                            font-size: 12px;
                            font-weight: 500;
                            text-transform: uppercase;
                        }

                        .code-copy-btn {
                            color: #61afef;
                            font-size: 12px;
                            cursor: pointer;
                            padding: 2px 8px;
                            border-radius: 4px;
                            transition: all 0.2s;

                            &:hover {
                                background-color: #3e4451;
                                color: #98c379;
                            }
                        }
                    }

                    .code-pre {
                        background-color: #282c34;
                        padding: 0;
                        margin: 0;

                        code {
                            background-color: transparent;
                            padding: 16px;
                            margin: 0;
                        }
                    }
                }

                blockquote {
                    border-left: 4px solid #ddd;
                    margin: 20px 0;
                    padding-left: 20px;
                    color: #666;
                    font-style: italic;
                }

                table {
                    width: 100%;
                    border-collapse: collapse;
                    margin: 20px 0;

                    th,
                    td {
                        border: 1px solid #ddd;
                        padding: 12px;
                        text-align: left;
                    }

                    th {
                        background: #f6f8fa;
                        font-weight: 600;
                    }
                }
            }
        }

        .article-footer {
            margin-top: 60px;
            padding-top: 30px;
            border-top: 1px solid #eee;

            .article-tags {
                margin-bottom: 20px;

                .el-tag {
                    margin-right: 10px;
                    margin-bottom: 10px;
                }
            }

            .article-actions {
                text-align: center;

                .action-buttons {
                    .el-button {
                        margin: 0 10px;
                    }
                }
            }
        }
    }
}

@media print {
    .preview-dialog {
        :deep(.el-dialog__header) {
            display: none;
        }
    }

    .preview-content {
        padding: 0 !important;

        .article-actions {
            display: none;
        }
    }
}
</style>