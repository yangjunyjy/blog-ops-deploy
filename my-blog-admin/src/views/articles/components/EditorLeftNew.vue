<template>
  <el-card shadow="never" class="editor-left">
    <!-- 标题输入 -->
    <div class="title-section">
      <el-input v-model="localTitle" placeholder="请输入文章标题" size="large" maxlength="100" show-word-limit
        class="title-input" />
    </div>

    <!-- 标签工具栏 -->
    <div class="tags-toolbar">
      <el-select v-model="localCategoryId" placeholder="选择分类" size="small" style="width: 150px">
        <el-option v-for="cat in categories" :key="cat.id" :label="cat.name" :value="cat.id" />
      </el-select>

      <el-select v-model="localTagIds" multiple placeholder="选择标签" size="small" style="width: 200px">
        <el-option v-for="tag in tags" :key="tag.id" :label="tag.name" :value="tag.id" />
      </el-select>

      <!-- 编辑器切换按钮 -->
      <el-button-group class="editor-switch" size="small">
        <el-button :type="editorType === 'markdown' ? 'primary' : 'default'" @click="switchEditor('markdown')">
          Markdown
        </el-button>
        <el-button :type="editorType === 'word' ? 'primary' : 'default'" @click="switchEditor('word')">
          Word
        </el-button>
      </el-button-group>
    </div>

    <!-- 工具栏和编辑器区域 -->
    <div class="editor-container">
      <!-- Word 编辑器 -->
      <template v-if="editorType === 'word'">
        <Toolbar class="editor-toolbar" :editor="editorRef" :defaultConfig="toolbarConfig" :mode="mode" />
        <Editor
          class="editor-content"
          key="word-editor"
          v-model="localContent"
          :defaultConfig="editorConfig"
          :mode="mode"
          @onCreated="handleEditorCreated"
          @onChange="handleEditorChange"
        />
      </template>

      <!-- Markdown 编辑器 -->
      <template v-else>
        <Toolbar class="editor-toolbar" :editor="editorRef" :defaultConfig="markdownToolbarConfig" :mode="mode" />
        <Editor
          class="editor-content"
          key="markdown-editor"
          v-model="localContent"
          :defaultConfig="markdownEditorConfig"
          :mode="mode"
          @onCreated="handleEditorCreated"
          @onChange="handleEditorChange"
        />
      </template>
    </div>

    <!-- 图片上传对话框 -->
    <el-dialog v-model="showImageUpload" title="上传图片" width="600px" :close-on-click-modal="false">
      <ImageUpload @upload-success="handleImageUploaded" @cancel="showImageUpload = false" />
    </el-dialog>
  </el-card>
</template>

<script setup>
import { ref, shallowRef, watch, onMounted, onBeforeUnmount } from 'vue'
import { ElMessage } from 'element-plus'
import '@wangeditor/editor/dist/css/style.css'
import { Editor, Toolbar } from '@wangeditor/editor-for-vue'
import ImageUpload from './ImageUpload.vue'
import { uploadImage as apiUploadImage } from '@/api/upload'

const props = defineProps({
  title: String,
  content: String,
  categories: Array,
  tags: Array,
  categoryId: [String, Number],
  tagIds: Array,
  editorType: {
    type: String,
    default: 'markdown'
  }
})

const emit = defineEmits([
  'update:title',
  'update:content',
  'update:categoryId',
  'update:tagIds',
  'content-change',
  'switch-editor',
  'editor-ready'
])

// 本地状态
const localTitle = ref(props.title)
const localContent = ref(props.content || '<p></p>')
const localCategoryId = ref(props.categoryId)
const localTagIds = ref(props.tagIds || [])
const showImageUpload = ref(false)

let lastSetContent = ''

// 修改监听器
watch(() => props.content, (val) => {
  console.log('EditorLeft收到content变化,长度:', val?.length || 0)

  // 如果内容相同，直接返回
  if (lastSetContent === val) {
    return
  }

  // 先更新本地内容，让 v-model 生效
  localContent.value = val || '<p></p>'

  // 如果编辑器已经创建，手动设置内容
  if (val && val.trim() && editorRef.value) {
    // 使用setTimeout异步执行,避免阻塞UI
    setTimeout(() => {
      try {
        console.log('开始设置编辑器HTML内容...')
        const startTime = Date.now()
        editorRef.value.setHtml(val)
        lastSetContent = val  // 更新标记
        console.log('编辑器内容设置完成,耗时:', Date.now() - startTime, 'ms')
      } catch (e) {
        console.error('更新编辑器内容失败:', e)
      }
    }, 100)
  } else {
    lastSetContent = val || '<p></p>'
  }
}, { immediate: true })

// 修改编辑器变化回调
const handleEditorChange = (editor) => {
  // 获取HTML内容
  const html = editor.getHtml()
  localContent.value = html

  // 只有当内容确实变化时才触发更新
  if (html !== lastSetContent) {
    lastSetContent = html  // 更新标记
    // 触发更新
    emit('update:content', html)
    emit('content-change', html)
  }
}

// 切换编辑器类型
const switchEditor = (type) => {
  emit('switch-editor', type)
}

// 监听props变化
watch(() => props.title, (val) => localTitle.value = val)
watch(() => props.categoryId, (val) => localCategoryId.value = val)
watch(() => props.tagIds, (val) => localTagIds.value = val || [])

// 监听编辑器类型变化，切换后确保内容被正确设置
watch(() => props.editorType, (newType, oldType) => {
  console.log('编辑器类型变化:', oldType, '->', newType)
  // 切换编辑器后，需要等待新编辑器创建完成
  // 这由组件的 key 属性自动处理
}, { immediate: false })

// 双向绑定
watch(localTitle, (val) => emit('update:title', val))
watch(localCategoryId, (val) => emit('update:categoryId', val))
watch(localTagIds, (val) => emit('update:tagIds', val))

// ==================== wangEditor 配置 ====================
const mode = 'default'
const editorRef = shallowRef()

// Word 工具栏配置
const toolbarConfig = {
  toolbarKeys: [
    // 字体样式
    'bold',
    'italic',
    'underline',
    'through',
    'color',
    'bgColor',
    'clearStyle',

    // 段落格式
    'fontSize',
    'fontFamily',
    'indent',
    'delIndent',
    'justifyLeft',
    'justifyCenter',
    'justifyRight',
    'justifyJustify',
    'lineHeight',

    // 列表
    'bulletedList',
    'numberedList',

    // 插入内容
    'insertLink',
    'editLink',
    'unLink',
    'codeBlock',
    'blockquote',
    'headerSelect',
    'header1',
    'header2',
    'header3',
    'header4',
    'header5',

    // 撤销重做
    'undo',
    'redo',

    // 媒体
    'insertImage',
    'deleteImage'
  ]
}

// Markdown 工具栏配置
const markdownToolbarConfig = {
  toolbarKeys: [
    'headerSelect',
    'bold',
    'italic',
    'through',
    'code',
    'bulletedList',
    'numberedList',
    'todo',
    'insertLink',
    'editLink',
    'unLink',
    'codeBlock',
    'undo',
    'redo',
    'insertImage',
    'deleteImage'
  ]
}

// Word 编辑器配置
const editorConfig = {
  placeholder: '开始写作吧...',
  scroll: true,
  maxLength: 50000,
  autoFocus: false,

  // 自定义快捷键
  customShortcut: {
    'mod+s': (editor) => {
      emit('content-change', localContent.value)
      ElMessage.info('已自动保存')
      return false
    },
    'mod+b': (editor) => {
      editor.execCommand('bold')
      return false
    },
    'mod+i': (editor) => {
      editor.execCommand('italic')
      return false
    },
    'mod+k': (editor) => {
      editor.execCommand('insertLink', { text: '链接文本', url: 'https://' })
      return false
    },
    'mod+shift+c': (editor) => {
      editor.execCommand('codeBlock')
      return false
    }
  },

  MENU_CONF: {
    uploadImage: {
      async customUpload(file, insertFn) {
        try {
          ElMessage.info('图片上传中...')

          if (file.size > 5 * 1024 * 1024) {
            ElMessage.error('图片大小不能超过5MB')
            return
          }

          const allowedTypes = ['image/jpeg', 'image/png', 'image/gif', 'image/webp', 'image/svg+xml']
          if (!allowedTypes.includes(file.type)) {
            ElMessage.error('只支持 JPG、PNG、GIF、WEBP、SVG 格式的图片')
            return
          }

          const formData = new FormData()
          formData.append('file', file)
          formData.append('type', 'article')

          const res = await apiUploadImage(formData)

          if (res.code === 200) {
            const imageUrl = res.data.url || res.data
            insertFn(imageUrl, file.name, imageUrl)
            ElMessage.success('图片上传成功')
          } else {
            ElMessage.error(res.msg || '图片上传失败')
          }
        } catch (error) {
          console.error('图片上传失败:', error)
          ElMessage.error('图片上传失败，请稍后重试')
        }
      }
    },
    insertLink: {
      checkLink(text, url) {
        if (!url) return false
        if (!url.startsWith('http://') && !url.startsWith('https://')) {
          url = 'https://' + url
        }
        return { text, url }
      }
    },
    codeSelectLang: {
      codeLangs: [
        { text: 'JavaScript', value: 'javascript' },
        { text: 'TypeScript', value: 'typescript' },
        { text: 'Python', value: 'python' },
        { text: 'Java', value: 'java' },
        { text: 'C++', value: 'cpp' },
        { text: 'C#', value: 'csharp' },
        { text: 'Go', value: 'go' },
        { text: 'Ruby', value: 'ruby' },
        { text: 'PHP', value: 'php' },
        { text: 'CSS', value: 'css' },
        { text: 'HTML', value: 'html' },
        { text: 'Shell', value: 'shell' },
        { text: 'SQL', value: 'sql' },
        { text: 'JSON', value: 'json' },
        { text: 'XML', value: 'xml' },
        { text: 'Markdown', value: 'markdown' },
        { text: 'Vue', value: 'vue' },
        { text: 'React', value: 'jsx' },
        { text: 'Dart', value: 'dart' },
        { text: 'Kotlin', value: 'kotlin' },
        { text: 'Swift', value: 'swift' },
        { text: 'Rust', value: 'rust' }
      ]
    }
  },

  fontSize: {
    fontSizeList: ['12px', '14px', '16px', '18px', '20px', '24px', '28px', '32px', '36px']
  },

  fontFamily: {
    fontFamilyList: ['微软雅黑', '宋体', '黑体', '楷体', 'Arial', 'Tahoma', 'Verdana', 'Times New Roman']
  },

  lineHeight: {
    lineHeightList: ['1', '1.5', '1.75', '2', '2.5', '3']
  }
}

// Markdown 编辑器配置
const markdownEditorConfig = {
  placeholder: '开始写作吧...',
  scroll: true,
  maxLength: 50000,
  autoFocus: false,

  // 自定义快捷键
  customShortcut: {
    'mod+s': (editor) => {
      emit('content-change', localContent.value)
      ElMessage.info('已自动保存')
      return false
    },
    'mod+b': (editor) => {
      editor.execCommand('bold')
      return false
    },
    'mod+i': (editor) => {
      editor.execCommand('italic')
      return false
    },
    'mod+k': (editor) => {
      editor.execCommand('insertLink', { text: '链接文本', url: 'https://' })
      return false
    },
    'mod+shift+c': (editor) => {
      editor.execCommand('codeBlock')
      return false
    }
  },

  MENU_CONF: {
    uploadImage: {
      async customUpload(file, insertFn) {
        try {
          ElMessage.info('图片上传中...')

          if (file.size > 5 * 1024 * 1024) {
            ElMessage.error('图片大小不能超过5MB')
            return
          }

          const allowedTypes = ['image/jpeg', 'image/png', 'image/gif', 'image/webp', 'image/svg+xml']
          if (!allowedTypes.includes(file.type)) {
            ElMessage.error('只支持 JPG、PNG、GIF、WEBP、SVG 格式的图片')
            return
          }

          const formData = new FormData()
          formData.append('file', file)
          formData.append('type', 'article')

          const res = await apiUploadImage(formData)

          if (res.code === 200) {
            const imageUrl = res.data.url || res.data
            insertFn(imageUrl, file.name, imageUrl)
            ElMessage.success('图片上传成功')
          } else {
            ElMessage.error(res.msg || '图片上传失败')
          }
        } catch (error) {
          console.error('图片上传失败:', error)
          ElMessage.error('图片上传失败，请稍后重试')
        }
      }
    },
    insertLink: {
      checkLink(text, url) {
        if (!url) return false
        if (!url.startsWith('http://') && !url.startsWith('https://')) {
          url = 'https://' + url
        }
        return { text, url }
      }
    },
    codeSelectLang: {
      codeLangs: [
        { text: 'JavaScript', value: 'javascript' },
        { text: 'TypeScript', value: 'typescript' },
        { text: 'Python', value: 'python' },
        { text: 'Java', value: 'java' },
        { text: 'C++', value: 'cpp' },
        { text: 'C#', value: 'csharp' },
        { text: 'Go', value: 'go' },
        { text: 'Ruby', value: 'ruby' },
        { text: 'PHP', value: 'php' },
        { text: 'CSS', value: 'css' },
        { text: 'HTML', value: 'html' },
        { text: 'Shell', value: 'shell' },
        { text: 'SQL', value: 'sql' },
        { text: 'JSON', value: 'json' },
        { text: 'XML', value: 'xml' },
        { text: 'Markdown', value: 'markdown' },
        { text: 'Vue', value: 'vue' },
        { text: 'React', value: 'jsx' },
        { text: 'Dart', value: 'dart' },
        { text: 'Kotlin', value: 'kotlin' },
        { text: 'Swift', value: 'swift' },
        { text: 'Rust', value: 'rust' }
      ]
    }
  }
}

// 编辑器创建完成回调
const handleEditorCreated = (editor) => {
  console.log('编辑器创建完成')
  editorRef.value = editor

  // 设置编辑器高度
  const editorDom = editor.getEditableContainer()
  if (editorDom) {
    editorDom.style.minHeight = '500px'
  }

  // 设置初始内容（优先使用 props.content）
  const contentToSet = props.content || localContent.value
  if (contentToSet && contentToSet.trim() && contentToSet !== '<p></p>') {
    try {
      console.log('编辑器创建完成，设置初始内容，长度:', contentToSet.length)
      editor.setHtml(contentToSet)
      localContent.value = contentToSet
      lastSetContent = contentToSet
    } catch (e) {
      console.error('设置编辑器内容失败:', e)
      editor.clear()
    }
  }

  // 通知父组件编辑器已就绪
  emit('editor-ready', true)
}

// 处理图片上传成功
const handleImageUploaded = (imageUrl) => {
  if (editorRef.value) {
    editorRef.value.insertImage({
      src: imageUrl,
      alt: '图片',
      title: '图片'
    })
  }
  showImageUpload.value = false
}

// 组件销毁时销毁编辑器
onBeforeUnmount(() => {
  if (editorRef.value == null) return
  editorRef.value.destroy()
})

// 暴露编辑器实例和插入图片方法，供父组件使用
defineExpose({
  insertImage: (imageUrl, alt = '图片') => {
    if (editorRef.value) {
      editorRef.value.insertImage({
        src: imageUrl,
        alt: alt,
        title: alt
      })
    }
  },
  insertText: (text) => {
    if (editorRef.value) {
      editorRef.value.insertText(text)
    }
  },
  getEditor: () => editorRef.value
})

// 自动保存功能
const autoSave = () => {
  if (localContent.value && localTitle.value) {
    setTimeout(() => {
      emit('content-change', localContent.value)
      autoSave()
    }, 30000)
  }
}

onMounted(() => {
  autoSave()
})
</script>

<style scoped lang="scss">
.editor-left {
  height: calc(100vh - 160px);
  display: flex;
  flex-direction: column;

  .title-section {
    margin-bottom: 20px;

    .title-input {
      :deep(.el-input__inner) {
        font-size: 24px;
        font-weight: 600;
        border: none;
        padding: 10px 0;

        &:focus {
          box-shadow: none;
        }
      }
    }
  }

  .tags-toolbar {
    display: flex;
    gap: 10px;
    margin-bottom: 20px;
    padding-bottom: 20px;
    border-bottom: 1px solid #f0f0f0;

    .editor-switch {
      margin-left: auto;
    }
  }

  .editor-container {
    flex: 1;
    display: flex;
    flex-direction: column;
    border: 1px solid #dcdfe6;
    border-radius: 6px;
    overflow: hidden;
    min-height: 600px;

    .editor-toolbar {
      border-bottom: 1px solid #dcdfe6;
      background-color: #f8f9fa;
      padding: 8px 12px;
      min-height: 44px;

      :deep(.w-e-bar) {
        background-color: transparent;
      }

      :deep(.w-e-bar-item) {
        margin-right: 4px;

        button {
          border-radius: 4px;

          &:hover {
            background-color: #e9ecef;
          }
        }
      }

      :deep(.w-e-bar-divider) {
        margin: 0 8px;
      }
    }

    .editor-content {
      flex: 1;
      overflow-y: auto;
      min-height: 500px;

      :deep(.w-e-text-container) {
        min-height: 500px;
        height: 100%;
        padding: 20px;
        line-height: 1.6;
        font-size: 16px;

        * {
          max-width: 100%;
        }

        [contenteditable='true'] {
          min-height: 500px;
        }
      }

      :deep(.w-e-scroll) {
        min-height: 500px !important;
      }

      // 标题样式
      h1,
      h2,
      h3,
      h4,
      h5,
      h6 {
        margin-top: 1.5em;
        margin-bottom: 0.5em;
        font-weight: 600;
        line-height: 1.25;
      }

      h1 {
        font-size: 2em;
      }

      h2 {
        font-size: 1.5em;
      }

      h3 {
        font-size: 1.25em;
      }

      // 段落样式
      p {
        margin: 0 0 1em 0;
        line-height: 1.8;
      }

      // CSDN 风格代码块
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
          padding: 16px;
          margin: 0;
          overflow-x: auto;

          code {
            background-color: transparent;
            padding: 0;
            border-radius: 0;
            font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', 'Consolas', monospace;
            font-size: 14px;
            line-height: 1.6;

            .hljs {
              background-color: transparent;
              color: #abb2bf;
            }

            .hljs-comment,
            .hljs-quote {
              color: #5c6370;
              font-style: italic;
            }

            .hljs-keyword,
            .hljs-selector-tag,
            .hljs-subst {
              color: #c678dd;
            }

            .hljs-number,
            .hljs-literal,
            .hljs-variable,
            .hljs-template-variable,
            .hljs-tag .hljs-attr {
              color: #d19a66;
            }

            .hljs-string,
            .hljs-doctag {
              color: #98c379;
            }

            .hljs-title,
            .hljs-section,
            .hljs-selector-id {
              color: #e06c75;
            }

            .hljs-type,
            .hljs-class .hljs-title {
              color: #e5c07b;
            }

            .hljs-tag,
            .hljs-name,
            .hljs-attribute {
              color: #61afef;
            }

            .hljs-regexp,
            .hljs-link {
              color: #56b6c2;
            }

            .hljs-symbol,
            .hljs-bullet {
              color: #61afef;
            }

            .hljs-built_in,
            .hljs-builtin-name {
              color: #e5c07b;
            }

            .hljs-meta {
              color: #61afef;
            }

            .hljs-deletion {
              background-color: #5c6370;
            }

            .hljs-addition {
              background-color: #98c379;
            }

            .hljs-emphasis {
              font-style: italic;
            }

            .hljs-strong {
              font-weight: bold;
            }
          }
        }
      }

      pre:not(.code-pre) {
        background-color: #f6f8fa;
        border-radius: 6px;
        padding: 16px;
        overflow: auto;
        margin: 1em 0;
        font-size: 14px;
        line-height: 1.45;

        code {
          background-color: transparent;
          padding: 0;
          border-radius: 0;
          font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, monospace;
        }
      }

      // 行内代码样式
      code:not(pre code):not(.hljs) {
        background-color: #f6f8fa;
        padding: 2px 6px;
        border-radius: 4px;
        font-size: 14px;
        font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, monospace;
        color: #e06c75;
      }

      // 引用样式
      blockquote {
        border-left: 4px solid #e5e7eb;
        margin: 1em 0;
        padding: 0 1em;
        color: #6b7280;

        p {
          margin: 0;
        }
      }

      // 列表样式
      ul,
      ol {
        padding-left: 2em;
        margin: 1em 0;

        li {
          margin-bottom: 0.5em;
        }
      }

      // 表格样式
      table {
        border-collapse: collapse;
        margin: 1em 0;
        width: 100%;

        th,
        td {
          border: 1px solid #e5e7eb;
          padding: 8px 12px;
          text-align: left;
        }

        th {
          background-color: #f9fafb;
          font-weight: 600;
        }
      }

      // 图片样式
      img {
        max-width: 100%;
        height: auto;
        border-radius: 6px;
        margin: 1em 0;
      }

      // 链接样式
      a {
        color: #3b82f6;
        text-decoration: none;

        &:hover {
          text-decoration: underline;
        }
      }
    }
  }
}

// 暗黑模式适配
html.dark {
  .editor-left {
    .tags-toolbar {
      border-bottom-color: #374151;
    }

    .editor-container {
      border-color: #374151;

      .editor-toolbar {
        background-color: #1f2937;
        border-bottom-color: #374151;

        :deep(.w-e-bar-item) {
          button {
            &:hover {
              background-color: #374151;
            }
          }
        }
      }

      .editor-content {
        :deep(.w-e-text-container) {
          background-color: #111827;
          color: #e5e7eb;

          pre {
            background-color: #1f2937;
          }

          code:not(pre code) {
            background-color: #1f2937;
            color: #e5e7eb;
          }

          blockquote {
            border-left-color: #4b5563;
            color: #9ca3af;
          }

          table {
            th,
            td {
              border-color: #374151;
            }

            th {
              background-color: #1f2937;
            }
          }
        }
      }
    }
  }
}

// 响应式设计
@media (max-width: 768px) {
  .editor-left {
    height: auto;
    min-height: calc(100vh - 160px);

    .tags-toolbar {
      flex-direction: column;
      align-items: flex-start;

      >* {
        width: 100%;
      }
    }

    .editor-container {
      .editor-toolbar {
        :deep(.w-e-bar-item) {
          margin-right: 2px;
        }

        :deep(.w-e-bar-divider) {
          margin: 0 4px;
        }
      }
    }
  }
}
</style>
