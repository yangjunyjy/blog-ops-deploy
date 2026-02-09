import { computed } from 'vue'
import { marked } from 'marked'
import hljs from 'highlight.js'
import 'highlight.js/styles/atom-one-dark.css'

export function useMarkdownRenderer() {
  // 配置 marked
  marked.setOptions({
    highlight: function(code, lang) {
      const language = hljs.getLanguage(lang) ? lang : 'plaintext'
      const highlighted = hljs.highlight(code, { language }).value
      return `<pre><code class="hljs language-${language}">${highlighted}</code></pre>`
    },
    langPrefix: 'hljs language-',
    breaks: true,
    gfm: true
  })

  // 渲染Markdown
  const renderMarkdown = (content) => {
    if (!content) return ''

    let html = marked.parse(content)

    // 后处理：为代码块添加 CSDN 风格包装
    html = html.replace(/<pre><code class="hljs language-(\w+)">([\s\S]*?)<\/code><\/pre>/g, (match, language, code) => {
      return `
        <div class="code-block-wrapper">
          <div class="code-header">
            <span class="code-language">${language}</span>
            <span class="code-copy-btn" onclick="copyCode(this)">复制</span>
          </div>
          <pre class="code-pre"><code class="hljs language-${language}">${code}</code></pre>
        </div>
      `
    })

    // 为行内代码添加样式
    html = html.replace(/<code>(?!.*hljs)(.*?)<\/code>/g, '<code class="inline-code">$1</code>')

    return html
  }

  return {
    renderMarkdown
  }
}

// 全局复制代码函数
if (typeof window !== 'undefined') {
  window.copyCode = function(btn) {
    const codeBlock = btn.closest('.code-block-wrapper').querySelector('code')
    const code = codeBlock.textContent

    navigator.clipboard.writeText(code).then(() => {
      btn.textContent = '已复制'
      setTimeout(() => {
        btn.textContent = '复制'
      }, 2000)
    }).catch(err => {
      console.error('复制失败:', err)
    })
  }
}
