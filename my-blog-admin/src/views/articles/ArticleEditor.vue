<template>
  <div class="article-editor">
    <EditorHeaderNew :publishing="publishing" :article-id="article.id" :editor-type="editorType" @back="handleBack"
      @save-draft="handleSaveDraft" @preview="handlePreview" @publish="handlePublish" @import="handleImport"
      @switch-editor="handleSwitchEditor" />

    <div class="editor-container">
      <el-row :gutter="20">
        <!-- 左侧编辑区 -->
        <el-col :xs="24" :sm="24" :md="16" :lg="16" :xl="16">
          <EditorLeftNew ref="editorLeftRef" v-model:title="article.title" v-model:content="article.content"
            :categories="categories" :tags="tags" :category-id="article.categoryId" :tag-ids="article.tags"
            :editor-type="editorType" @update:category-id="val => article.categoryId = val"
            @update:tag-ids="val => article.tags = val" @content-change="handleContentChange"
            @switch-editor="handleSwitchEditor" @editor-ready="handleEditorReady" />
        </el-col>

        <!-- 右侧设置区 -->
        <el-col :xs="24" :sm="24" :md="8" :lg="8" :xl="8">
          <EditorRight v-model:article="article" :word-count="wordCount" :image-count="imageCount"
            :reading-time="readingTime" @cover-uploaded="handleCoverUploaded" @update:article="val => article = val" />
        </el-col>
      </el-row>
    </div>

    <!-- 预览对话框 -->
    <PreviewDialog v-model="previewVisible" :article="article" :rendered-content="renderedContent" />
  </div>
</template>

<script setup>
defineOptions({
  name: 'ArticleEditor'
})
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import EditorHeaderNew from './components/EditorHeaderNew.vue'
import EditorLeftNew from './components/EditorLeftNew.vue'
import EditorRight from './components/EditorRight.vue'
import PreviewDialog from './components/PreviewDialog.vue'
import { getArticleDetail, createArticle, updateArticle } from '@/api/article'
import { getCategoryList } from '@/api/category'
import { getTagList } from '@/api/tag'
import { marked } from 'marked'

const route = useRoute()
const router = useRouter()

// 从 localStorage 获取用户信息
const getUserInfo = () => {
  try {
    return JSON.parse(localStorage.getItem('userInfo') || '{}')
  } catch {
    return {}
  }
}

// 状态
const publishing = ref(false)
const previewVisible = ref(false)
const categories = ref([])
const tags = ref([])
const editorType = ref('markdown') // 'markdown' | 'word'
const editorLeftRef = ref(null) // 编辑器组件引用
const isEditorReady = ref(false) // 编辑器是否就绪
const pendingImages = ref([]) // 待插入的图片列表

// 文章数据
const userInfo = getUserInfo()
const article = reactive({
  id: null,
  title: '',
  content: '',
  cover: '',
  slug: '',
  categoryId: null,
  tags: [],
  summary: '',
  keywords: '',
  is_top: false,
  is_recommended: false,
  allow_comment: true,
  status: 0,
  publishedAt: null,
  author_id: userInfo.id || 0
})

// 统计信息
const wordCount = computed(() => {
  // 编辑器中始终是 HTML，直接统计纯文本
  const plainText = article.content.replace(/<[^>]*>/g, '')
  return plainText.replace(/[^\u4e00-\u9fa5\w]/g, '').length
})

const imageCount = computed(() => {
  const matches = article.content.match(/<img[^>]+src="[^"]+"/gi)
  return matches ? matches.length : 0
})

const readingTime = computed(() => {
  const minutes = Math.ceil(wordCount.value / 300)
  return `${minutes} 分钟`
})

// 计算渲染后的内容 - 仅用于预览
const renderedContent = computed(() => {
  // 编辑器中始终是 HTML 格式（所见即所得）
  return article.content || ''
})

// 加载数据
const loadArticle = async (id) => {
  try {
    const res = await getArticleDetail(id)
    const data = res.data

    // 如果内容是 Markdown，转换为 HTML 显示
    if (data.content && !/<[a-z][\s\S]*>/i.test(data.content)) {
      data.content = markdownToHtml(data.content)
    }

    Object.assign(article, data)
  } catch (error) {
    console.error('加载文章失败:', error)
    ElMessage.error('加载文章失败')
  }
}

const loadCategories = async () => {
  try {
    const res = await getCategoryList({ page: 1, page_size: 100 })
    categories.value = res.data.items || []
  } catch (error) {
    console.error('加载分类失败:', error)
  }
}

const loadTags = async () => {
  try {
    const res = await getTagList({ page: 1, page_size: 100 })
    tags.value = res.data.items || []
  } catch (error) {
    console.error('加载标签失败:', error)
  }
}

// 事件处理
const handleBack = () => {
  router.back()
}

const handleCoverUploaded = (url) => {
  article.cover = url
}

const handleContentChange = (content) => {
  article.content = content
}

const handleEditorReady = () => {
  console.log('编辑器已就绪')
  isEditorReady.value = true

  // 如果有待插入的图片，现在插入
  if (pendingImages.value.length > 0) {
    insertPendingImages()
  }
}

const insertPendingImages = () => {
  if (!editorLeftRef.value || pendingImages.value.length === 0) return

  console.log('插入待处理的图片，数量:', pendingImages.value.length)
  let insertCount = 0

  pendingImages.value.forEach((url) => {
    try {
      editorLeftRef.value.insertImage(url, '文档图片')
      insertCount++
    } catch (e) {
      console.error('插入图片失败:', e)
    }
  })

  console.log('图片插入完成，成功:', insertCount, '张')
  pendingImages.value = []
}

const handleSwitchEditor = (type) => {
  editorType.value = type
  ElMessage.success(`已切换到 ${type === 'word' ? 'Word' : 'Markdown'} 编辑器`)
}

// 处理导入的文章
const handleImport = (data) => {
  console.log('收到导入事件:', data)
  console.log('当前文章状态:', { title: article.title, hasContent: !!article.content })

  // 如果当前文章已有内容，提示用户确认
  if (article.title || (article.content && article.content !== '<p></p>')) {
    console.log('文章已有内容,显示确认对话框')
    ElMessageBox.confirm(
      '当前文章已有内容，导入将会覆盖，是否继续？',
      '提示',
      {
        confirmButtonText: '继续导入',
        cancelButtonText: '取消',
        type: 'warning'
      }
    ).then(() => {
      console.log('用户确认导入')
      applyImport(data)
    }).catch((action) => {
      console.log('用户取消导入:', action)
      // action可能是'cancel'或'close'
    })
  } else {
    console.log('直接导入，无需确认')
    applyImport(data)
  }
}

// 将 Markdown 转换为 HTML（用于显示）
const markdownToHtml = (markdown) => {
  if (!markdown) return ''
  try {
    // 检查内容大小,如果太大则直接返回原内容(避免卡死)
    if (markdown.length > 100000) {
      console.warn('Markdown内容太大,直接显示原内容')
      return `<pre>${escapeHtml(markdown)}</pre>`
    }
    return marked.parse(markdown)
  } catch (e) {
    console.error('Markdown 解析错误:', e)
    return `<pre>${escapeHtml(markdown)}</pre>`
  }
}

// HTML转义函数
const escapeHtml = (text) => {
  if (!text) return ''
  return text
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#039;')
}

// 将 HTML 转换为 Markdown（用于保存）
const htmlToMarkdown = (html) => {
  if (!html) return ''
  // 简单的 HTML 到 Markdown 转换
  return html
    .replace(/<h1>(.*?)<\/h1>/gi, '\n# $1\n')
    .replace(/<h2>(.*?)<\/h2>/gi, '\n## $1\n')
    .replace(/<h3>(.*?)<\/h3>/gi, '\n### $1\n')
    .replace(/<h4>(.*?)<\/h4>/gi, '\n#### $1\n')
    .replace(/<h5>(.*?)<\/h5>/gi, '\n##### $1\n')
    .replace(/<h6>(.*?)<\/h6>/gi, '\n###### $1\n')
    .replace(/<strong>(.*?)<\/strong>/gi, '**$1**')
    .replace(/<b>(.*?)<\/b>/gi, '**$1**')
    .replace(/<em>(.*?)<\/em>/gi, '*$1*')
    .replace(/<i>(.*?)<\/i>/gi, '*$1*')
    .replace(/<code\s+class="language-(\w+)">(.*?)<\/code>/gis, '```$1\n$2\n```')
    .replace(/<code>(.*?)<\/code>/gi, '`$1`')
    .replace(/<pre><code>(.*?)<\/code><\/pre>/gis, '```\n$1\n```')
    .replace(/<pre\s+[^>]*><code\s+class="language-(\w+)">(.*?)<\/code><\/pre>/gis, '```$1\n$2\n```')
    .replace(/<p>(.*?)<\/p>/gi, '\n\n$1\n\n')
    .replace(/<br\s*\/?>/gi, '\n')
    .replace(/<a\s+href="([^"]+)">(.*?)<\/a>/gi, '[$2]($1)')
    .replace(/<img\s+src="([^"]+)"[^>]*alt="([^"]*)"[^>]*>/gi, '![$2]($1)')
    .replace(/<img\s+src="([^"]+)"[^>]*>/gi, '![]($1)')
    .replace(/<blockquote>(.*?)<\/blockquote>/gis, '\n> $1\n')
    .replace(/<ul>([\s\S]*?)<\/ul>/gis, (_, content) => {
      const items = content.replace(/<li>(.*?)<\/li>/gi, '\n- $1')
      return items + '\n'
    })
    .replace(/<ol>([\s\S]*?)<\/ol>/gis, (_, content) => {
      let index = 1
      const items = content.replace(/<li>(.*?)<\/li>/gi, () => `\n${index++}. $1`)
      return items + '\n'
    })
    .replace(/<hr\s*\/?>/gi, '\n---\n')
    .replace(/&lt;/g, '<')
    .replace(/&gt;/g, '>')
    .replace(/&amp;/g, '&')
    .replace(/&nbsp;/g, ' ')
    .replace(/<\/?[^>]+(>|$)/g, '') // 移除剩余的 HTML 标签
    .replace(/\n{3,}/g, '\n\n') // 合并多余的换行
    .trim()
}

// 应用导入的内容
const applyImport = async (data) => {
  try {
    console.log('开始应用导入数据:', data)
    console.log('导入内容类型:', data.contentType || 'unknown')

    // 直接设置标题
    article.title = data.title

    // 根据内容类型处理
    if (data.contentType === 'markdown' || data._isMarkdown) {
      // 导入的是 Markdown 源码
      // 将 Markdown 渲染为 HTML 显示在编辑器中
      const htmlContent = markdownToHtml(data.content)
      article.content = htmlContent
      // 保存原始 Markdown 源码到特殊字段
      article._markdownSource = data.content
      console.log('Markdown 内容渲染为 HTML，HTML 长度:', htmlContent.length)
      console.log('原始 Markdown 长度:', data.content.length)
    } else {
      // 导入的是 HTML（已渲染的内容）
      article.content = data.content || ''
      article._markdownSource = data.content
      console.log('设置的 HTML 内容，长度:', article.content.length)
    }

    // 设置其他字段
    article.summary = data.summary || ''

    ElMessage.success('导入成功')
    console.log('导入应用完成')
  } catch (error) {
    console.error('导入失败:', error)
    ElMessage.error('导入失败: ' + (error.message || error))
  }
}



const handleSaveDraft = async () => {
  try {
    article.status = 0
    // 如果没有摘要，自动生成（支持 HTML 和 Markdown）
    if (!article.summary || article.summary.trim() === '') {
      const isHtml = /<[a-z][\s\S]*>/i.test(article.content)
      const plainText = isHtml ? article.content.replace(/<[^>]*>/g, '') : article.content
      const summary = plainText.substring(0, 200)
      article.summary = summary + (summary.length >= 200 ? '...' : '')
    }
    // 如果没有slug，自动生成
    if (!article.slug || article.slug.trim() === '') {
      article.slug = article.title.toLowerCase().replace(/[^a-z0-9\u4e00-\u9fa5]+/g, '-').replace(/-+/g, '-').replace(/^-+|-+$/g, '')
    }
    // 优先使用原始 Markdown 源码，如果没有则将编辑器中的 HTML 转换为 Markdown
    let contentToSave = article._markdownSource || htmlToMarkdown(article.content)
    const articleToSave = {
      title: article.title,
      content: contentToSave,
      summary: article.summary || '未提供摘要',
      cover: article.cover,
      categoryId: article.category_id,
      tags: article.tag_ids || [],
      status: article.status,
      slug: article.slug
    }
    if (article.id) {
      await updateArticle(article.id, articleToSave)
    } else {
      const res = await createArticle(articleToSave)
      article.id = res.data.id
    }
    ElMessage.success('草稿保存成功')
  } catch (error) {
    console.error('保存草稿失败:', error)
    ElMessage.error('保存失败: ' + (error.message || error.msg || '未知错误'))
  }
}

const handlePreview = () => {
  if (!article.title) {
    ElMessage.warning('请输入文章标题')
    return
  }
  previewVisible.value = true
}

const handlePublish = async () => {
  if (!article.title) {
    ElMessage.warning('请输入文章标题')
    return
  }
  if (!article.content) {
    ElMessage.warning('请输入文章内容')
    return
  }

  // 如果没有摘要，自动生成
  if (!article.summary || article.summary.trim() === '') {
    const isHtml = /<[a-z][\s\S]*>/i.test(article.content)
    const plainText = isHtml ? article.content.replace(/<[^>]*>/g, '') : article.content
    const summary = plainText.substring(0, 200)
    article.summary = summary + (summary.length >= 200 ? '...' : '')
  }

  // 如果没有slug，自动生成
  if (!article.slug || article.slug.trim() === '') {
    article.slug = article.title.toLowerCase().replace(/[^a-z0-9\u4e00-\u9fa5]+/g, '-').replace(/-+/g, '-').replace(/^-+|-+$/g, '')
  }

  // 优先使用原始 Markdown 源码，如果没有则将编辑器中的 HTML 转换为 Markdown
  let contentToSave = article._markdownSource || htmlToMarkdown(article.content)

  const articleToSave = {
    title: article.title,
    content: contentToSave,
    summary: article.summary || '未提供摘要',
    cover: article.cover,
    categoryId: article.category_id,
    tags: article.tag_ids || [],
    status: 1,
    is_top: article.is_top || false,
    slug: article.slug
  }

  publishing.value = true
  try {
    if (article.id) {
      await updateArticle(article.id, articleToSave)
    } else {
      const res = await createArticle(articleToSave)
      article.id = res.data.id
    }
    ElMessage.success('发布成功')
    router.push('/articles')
  } catch (error) {
    console.error('发布失败:', error)
    ElMessage.error('发布失败: ' + (error.message || error.msg || '未知错误'))
  } finally {
    publishing.value = false
  }
}

// 初始化
onMounted(() => {
  loadCategories()
  loadTags()

  const id = route.params.id
  if (id) {
    loadArticle(id)
  }
})
</script>

<style scoped lang="scss">
.article-editor {
  min-height: 100vh;
  background-color: #f0f2f5;
  display: flex;
  flex-direction: column;

  .editor-container {
    flex: 1;
    padding: 20px;
    padding-top: 10px;
  }
}
</style>