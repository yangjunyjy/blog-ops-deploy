import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import { marked, Renderer } from 'marked'
import hljs from 'highlight.js'
import 'highlight.js/styles/atom-one-dark.css'
import 'dayjs/locale/zh-cn'

dayjs.extend(relativeTime)
dayjs.locale('zh-cn')

// 自定义渲染器，集成 highlight.js
const renderer = new Renderer()
renderer.code = function(code, infoString) {
  // marked.js 5.x+ 传递的是对象 { text, lang, escaped }
  let codeText = ''
  let lang = ''

  if (typeof code === 'object' && code.text) {
    // marked.js 5.x+ 格式
    codeText = code.text
    lang = code.lang || ''
  } else if (typeof code === 'string') {
    // 旧版本格式（兼容）
    codeText = code
    lang = infoString || ''
  }

  let highlighted = codeText

  try {
    if (lang && hljs.getLanguage(lang)) {
      const result = hljs.highlight(codeText, { language: lang })
      highlighted = result.value || codeText
    } else {
      const result = hljs.highlightAuto(codeText)
      highlighted = result.value || codeText
    }
  } catch (e) {
    console.error('highlight.js error:', e)
    highlighted = codeText
  }

  // 不需要转义 HTML，因为 highlight.js 已经处理了
  // 生成行号
  const lines = highlighted.split('\n')
  const lineNumbers = lines.map((_, index) => `<span class="code-line-number">${index + 1}</span>`).join('')

  // 将代码按行包装，用于对齐
  const codeLines = lines.map((line, index) => {
    if (index === lines.length - 1 && line === '') return ''
    return `<div class="code-line">${line || ' '}</div>`
  }).join('')

  const displayLanguage = lang || 'plaintext'

  return `
    <div class="code-block-wrapper" data-language="${displayLanguage}">
      <div class="code-header">
        <span class="code-language">${displayLanguage}</span>
        <button class="code-copy-btn" onclick="copyCode(this)">
          <svg viewBox="0 0 24 24" width="14" height="14">
            <path fill="currentColor" d="M16 1H4c-1.1 0-2 .9-2 2v14h2V3h12V1zm3 4H8c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h11c1.1 0 2-.9 2-2V7c0-1.1-.9-2-2-2zm0 16H8V7h11v14z"/>
          </svg>
          复制
        </button>
      </div>
      <div class="code-body">
        <div class="code-line-numbers">${lineNumbers}</div>
        <div class="code-content">
          <code class="hljs language-${displayLanguage}">${highlighted}</code>
        </div>
      </div>
    </div>
  `
}

// 配置 marked - 兼容 marked 17.x
marked.use({
  renderer: renderer,
  breaks: true,
  gfm: true
})

// Markdown 转 HTML（智能判断：如果是HTML则直接返回，如果是Markdown则转换）
export const markdownToHtml = (content) => {
  if (!content) return ''

  // 判断是否已经是 HTML（包含常见的 HTML 标签）
  const isHtml = /<[a-z][\s\S]*>/i.test(content)

  // 如果已经是 HTML，直接返回（wangEditor 生成的）
  if (isHtml) {
    return content
  }

  // 否则当作 Markdown 处理
  try {
    return marked.parse(content)
  } catch (e) {
    console.error('marked.js parse error:', e)
    return content
  }
}

// 格式化日期
export const formatDate = (date, format = 'YYYY-MM-DD') => {
  return dayjs(date).format(format)
}

// 格式化相对时间
export const formatRelativeTime = (date) => {
  return dayjs(date).fromNow()
}

// 截取文本
export const truncateText = (text, length = 150) => {
  if (!text) return ''
  if (text.length <= length) return text
  return text.substring(0, length) + '...'
}
