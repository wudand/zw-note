<template>
  <div class="markdown-editor">
    <!-- 工具栏 -->
    <div class="markdown-editor__toolbar">
        <div class="markdown-editor__toolbar-left">
            <!-- 撤销/重做组 -->
            <div class="toolbar-group">
                <button
                    type="button"
                    class="toolbar-btn"
                    :disabled="!canUndo"
                    @click="handleUndo"
                    title="撤销"
                >
                <svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor">
                    <path d="M12.5 8c-2.65 0-5.05.99-6.9 2.6L2 7v9h9l-3.62-3.62c1.39-1.16 3.16-1.88 5.12-1.88 3.54 0 6.55 2.31 7.6 5.5l2.37-.78C21.08 11.03 17.15 8 12.5 8z"/>
                </svg>
                </button>
                <button
                    type="button"
                    class="toolbar-btn"
                    :disabled="!canRedo"
                    @click="handleRedo"
                    title="重做"
                >
                <svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor">
                    <path d="M18.4 10.6C16.55 8.99 14.15 8 11.5 8c-4.65 0-8.58 3.03-9.96 7.22L3.9 16c1.05-3.19 4.05-5.5 7.6-5.5 1.95 0 3.73.72 5.12 1.88L13 19h9v-9l-3.6 3.6z"/>
                </svg>
                </button>
            </div>

            <div class="toolbar-divider"></div>

            <!-- 标题组 -->
            <div class="toolbar-group">
                <button
                    type="button"
                    class="toolbar-btn"
                    @click="insertHeading(1)"
                    title="标题 1"
                >
                H1
                </button>
                <button
                    type="button"
                    class="toolbar-btn"
                    @click="insertHeading(2)"
                    title="标题 2"
                >
                H2
                </button>
                <button
                    type="button"
                    class="toolbar-btn"
                    @click="insertHeading(3)"
                    title="标题 3"
                >
                H3
                </button>
                <button
                    type="button"
                    class="toolbar-btn"
                    @click="insertHeading(4)"
                    title="标题 4"
                >
                H4
                </button>
                <button
                    type="button"
                    class="toolbar-btn"
                    @click="insertHeading(5)"
                    title="标题 5"
                >
                H5
                </button>
                <button
                    type="button"
                    class="toolbar-btn"
                    @click="insertHeading(6)"
                    title="标题 6"
                >
                H6
                </button>
            </div>

            <div class="toolbar-divider"></div>

            <!-- 文本格式化组 -->
            <div class="toolbar-group">
                <button
                    type="button"
                    class="toolbar-btn"
                    @click="wrapText('**', '**')"
                    title="粗体"
                >
                <strong>B</strong>
                </button>
                <button
                    type="button"
                    class="toolbar-btn"
                    @click="wrapText('*', '*')"
                    title="斜体"
                >
                <em>I</em>
                </button>
                <button
                    type="button"
                    class="toolbar-btn"
                    @click="wrapText('~~', '~~')"
                    title="删除线"
                >
                <span style="text-decoration: line-through">S</span>
                </button>
                <button
                    type="button"
                    class="toolbar-btn"
                    @click="wrapText('`', '`')"
                    title="行内代码"
                >
                &lt;/&gt;
                </button>
            </div>

            <div class="toolbar-divider"></div>

            <!-- 文本对齐组 -->
            <div class="toolbar-group">
                <button
                    type="button"
                    class="toolbar-btn"
                    @click="insertAlign('left')"
                    title="左对齐"
                >
                <svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor">
                    <path d="M3 3h18v2H3V3zm0 4h12v2H3V7zm0 4h18v2H3v-2zm0 4h12v2H3v-2zm0 4h18v2H3v-2z"/>
                </svg>
                </button>
                <button
                    type="button"
                    class="toolbar-btn"
                    @click="insertAlign('center')"
                    title="居中"
                >
                <svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor">
                    <path d="M3 3h18v2H3V3zm3 4h12v2H6V7zm-3 4h18v2H3v-2zm3 4h12v2H6v-2zm-3 4h18v2H3v-2z"/>
                </svg>
                </button>
                <button
                    type="button"
                    class="toolbar-btn"
                    @click="insertAlign('right')"
                    title="右对齐"
                >
                <svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor">
                    <path d="M3 3h18v2H3V3zm6 4h12v2H9V7zm-6 4h18v2H3v-2zm6 4h12v2H9v-2zm-6 4h18v2H3v-2z"/>
                </svg>
                </button>
                <button
                    type="button"
                    class="toolbar-btn"
                    @click="insertAlign('justify')"
                    title="两端对齐"
                >
                <svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor">
                    <path d="M3 3h18v2H3V3zm0 4h18v2H3V7zm0 4h18v2H3v-2zm0 4h18v2H3v-2zm0 4h18v2H3v-2z"/>
                </svg>
                </button>
            </div>

            <div class="toolbar-divider"></div>

            <!-- 列表组 -->
            <div class="toolbar-group">
                <button
                    type="button"
                    class="toolbar-btn"
                    @click="insertList('unordered')"
                    title="无序列表"
                >
                <svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor">
                    <circle cx="5" cy="6" r="1.5"/>
                    <circle cx="5" cy="12" r="1.5"/>
                    <circle cx="5" cy="18" r="1.5"/>
                    <path d="M10 6h11M10 12h11M10 18h11" stroke="currentColor" stroke-width="2"/>
                </svg>
                </button>
                <button
                    type="button"
                    class="toolbar-btn"
                    @click="insertList('ordered')"
                    title="有序列表"
                >
                <svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor">
                    <path d="M3 6h1v1H3V6zm0 6h1v1H3v-1zm0 6h1v1H3v-1zm3-12h15v2H6V6zm0 6h15v2H6v-2zm0 6h15v2H6v-2z"/>
                </svg>
                </button>
                <button
                    type="button"
                    class="toolbar-btn"
                    @click="insertList('task')"
                    title="任务列表"
                >
                <svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor">
                    <path d="M3 5h1v1H3V5zm0 6h1v1H3v-1zm0 6h1v1H3v-1zm3-10h15v2H6V1zm0 6h15v2H6V7zm0 6h15v2H6v-2z"/>
                    <path d="M5 5l2 2 4-4" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round"/>
                </svg>
                </button>
            </div>

            <div class="toolbar-divider"></div>

            <!-- 插入组 -->
            <div class="toolbar-group">
                <button
                    type="button"
                    class="toolbar-btn"
                    @click="insertLink"
                    title="插入链接"
                >
                <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71"/>
                    <path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71"/>
                </svg>
                </button>
                <button
                    type="button"
                    class="toolbar-btn"
                    @click="insertImage"
                    title="插入图片"
                >
                <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2">
                    <rect x="3" y="3" width="18" height="18" rx="2" ry="2"/>
                    <circle cx="8.5" cy="8.5" r="1.5"/>
                    <path d="M21 15l-5-5L5 21"/>
                </svg>
                </button>
                <button
                    type="button"
                    class="toolbar-btn"
                    @click="insertCodeBlock"
                    title="插入代码块"
                >
                <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M16 18l6-6-6-6M8 6l-6 6 6 6"/>
                </svg>
                </button>
                <button
                    type="button"
                    class="toolbar-btn"
                    @click="insertHorizontalRule"
                    title="插入水平线"
                >
                <svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor">
                    <path d="M3 12h18v2H3v-2z"/>
                </svg>
                </button>
                <button
                    type="button"
                    class="toolbar-btn"
                    @click="insertBlockquote"
                    title="插入引用"
                >
                <svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor">
                    <path d="M6 17h3l2-4V7H5v6h3zm8 0h3l2-4V7h-6v6h3z"/>
                </svg>
                </button>
                <button
                    type="button"
                    class="toolbar-btn"
                    @click="insertTable"
                    title="插入表格"
                >
                <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M3 3h18v18H3V3zM3 9h18M9 3v18"/>
                </svg>
                </button>
            </div>            
        </div>
        <div class="markdown-editor__toolbar-right">
            <!-- 保存 -->
            <button
                type="button"
                class="toolbar-btn toolbar-btn--save"
                @click="save"
                title="保存"
            >
                <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
                    <polyline points="7 10 12 15 17 10"/>
                    <line x1="12" y1="15" x2="12" y2="3"/>
                </svg>
            </button>

            <!-- 编辑/编辑和预览/全屏预览 -->
            <div class="toolbar-group">
                <button
                    type="button"
                    class="toolbar-btn toolbar-action-btn"
                    :class="{ 'toolbar-action-btn--active': mode === 'edit' }"
                    @click="switchMode('edit')"
                    title="编辑"
                >
                    <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                        <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                    </svg>
                </button>
                <button
                    type="button"
                    class="toolbar-btn toolbar-action-btn"
                    :class="{ 'toolbar-action-btn--active': mode === 'preview' }"
                    @click="switchMode('preview')"
                    title="编辑和预览"
                >
                    <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <rect x="2" y="3" width="20" height="14" rx="2" ry="2"/>
                        <line x1="8" y1="21" x2="16" y2="21"/>
                        <line x1="12" y1="17" x2="12" y2="21"/>
                        <line x1="7" y1="8" x2="7" y2="8.01"/>
                        <line x1="12" y1="8" x2="12" y2="8.01"/>
                        <line x1="17" y1="8" x2="17" y2="8.01"/>
                    </svg>
                </button>
                <button
                    type="button"
                    class="toolbar-btn toolbar-action-btn"
                    :class="{ 'toolbar-action-btn--active': mode === 'fullscreen' }"
                    @click="switchMode('fullscreen')"
                    title="全屏预览"
                >
                    <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <path d="M8 3H5a2 2 0 0 0-2 2v3m18 0V5a2 2 0 0 0-2-2h-3m0 18h3a2 2 0 0 0 2-2v-3M3 16v3a2 2 0 0 0 2 2h3"/>
                    </svg>
                </button>
            </div>
        </div>
    </div>

    <!-- 编辑区域 -->
    <div class="markdown-editor__editor-wrapper">
      <!-- 编辑模式：只显示编辑区域 -->
      <template v-if="mode === 'edit'">
        <textarea
          ref="editorRef"
          v-model="content"
          class="markdown-editor__editor"
          :placeholder="placeholder"
          @input="handleInput"
          @keydown="handleKeydown"
        ></textarea>
      </template>

      <!-- 编辑和预览模式：左右分屏 -->
      <template v-else-if="mode === 'preview'">
        <Splitpanes class="markdown-editor__splitpanes">
          <Pane :size="50" class="markdown-editor__pane">
            <textarea
              ref="editorRef"
              v-model="content"
              class="markdown-editor__editor"
              :placeholder="placeholder"
              @input="handleInput"
              @keydown="handleKeydown"
            ></textarea>
          </Pane>
          <Pane :size="50" class="markdown-editor__pane">
            <div class="markdown-editor__preview" v-html="renderedContent"></div>
          </Pane>
        </Splitpanes>
      </template>

      <!-- 全屏预览模式：只显示预览区域 -->
      <template v-else-if="mode === 'fullscreen'">
        <div class="markdown-editor__preview markdown-editor__preview--fullscreen" v-html="renderedContent"></div>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, nextTick, computed } from 'vue'

interface Props {
  modelValue?: string
  placeholder?: string
  height?: string
}

interface Emits {
  (e: 'update:modelValue', value: string): void
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: '',
  placeholder: '请输入内容...',
  height: '500px',
})

const emit = defineEmits<Emits>()

const editorRef = ref<HTMLTextAreaElement>()
const content = ref(props.modelValue)
const history = ref<string[]>([props.modelValue])
const historyIndex = ref(0)
const canUndo = ref(false)
const canRedo = ref(false)

// 编辑模式
const mode = ref<'edit' | 'preview' | 'fullscreen'>('preview')

// 切换模式
function switchMode(newMode: 'edit' | 'preview' | 'fullscreen') {
  mode.value = newMode
}

// 保存功能
function save() {
  // TODO: 实现保存逻辑
  console.log('保存文档:', content.value)
  emit('update:modelValue', content.value)
}

// Markdown 渲染函数
function renderMarkdown(markdown: string): string {
  if (!markdown) return ''
  
  let html = markdown
  
  // 先处理代码块（避免代码块内的内容被其他规则处理）
  const codeBlocks: string[] = []
  html = html.replace(/```[\s\S]*?```/g, (match) => {
    const index = codeBlocks.length
    codeBlocks.push(match)
    // 使用大括号占位符，避免与 Markdown 的下划线语法冲突
    return `{{MARKDOWN-CODE-BLOCK-${index}}}`
  })
  
  // 转义 HTML 特殊字符（除了我们允许的标签）
  html = html
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
  
  // 恢复允许的 HTML 标签（如对齐标签）
  html = html
    .replace(/&lt;div align="(left|center|right|justify)"&gt;/g, '<div align="$1">')
    .replace(/&lt;\/div&gt;/g, '</div>')
  
  // 标题（按从大到小顺序处理，避免冲突）
  html = html.replace(/^###### (.*$)/gim, '<h6>$1</h6>')
  html = html.replace(/^##### (.*$)/gim, '<h5>$1</h5>')
  html = html.replace(/^#### (.*$)/gim, '<h4>$1</h4>')
  html = html.replace(/^### (.*$)/gim, '<h3>$1</h3>')
  html = html.replace(/^## (.*$)/gim, '<h2>$1</h2>')
  html = html.replace(/^# (.*$)/gim, '<h1>$1</h1>')
  
  // 水平线
  html = html.replace(/^---$/gim, '<hr>')
  html = html.replace(/^\*\*\*$/gim, '<hr>')
  html = html.replace(/^___$/gim, '<hr>')
  
  // 行内代码
  html = html.replace(/`([^`]+)`/g, '<code>$1</code>')
  
  // 粗体
  html = html.replace(/\*\*([^*]+)\*\*/g, '<strong>$1</strong>')
  html = html.replace(/__([^_]+)__/g, '<strong>$1</strong>')
  
  // 斜体（在粗体之后处理，避免冲突）
  html = html.replace(/(?<!\*)\*([^*]+)\*(?!\*)/g, '<em>$1</em>')
  html = html.replace(/(?<!_)_([^_]+)_(?!_)/g, '<em>$1</em>')
  
  // 删除线
  html = html.replace(/~~([^~]+)~~/g, '<del>$1</del>')
  
  // 链接
  html = html.replace(/\[([^\]]+)\]\(([^)]+)\)/g, '<a href="$2" target="_blank" rel="noopener noreferrer">$1</a>')
  
  // 图片
  html = html.replace(/!\[([^\]]*)\]\(([^)]+)\)/g, '<img src="$2" alt="$1">')
  
  // 恢复代码块
  codeBlocks.forEach((codeBlock, index) => {
    const code = codeBlock.replace(/```[\w]*\n?/g, '').replace(/```/g, '').trim()
    html = html.replace(`{{MARKDOWN-CODE-BLOCK-${index}}}`, `<pre><code>${escapeHtml(code)}</code></pre>`)
  })
  
  // 按行处理
  const lines = html.split('\n')
  const result: string[] = []
  let inList = false
  let listType: 'ul' | 'ol' | 'task' | null = null
  
  for (let i = 0; i < lines.length; i++) {
    const line = lines[i]
    if (line === undefined) continue
    const trimmed = line.trim()
    
    // 空行
    if (!trimmed) {
      if (inList) {
        result.push(`</${listType === 'ol' ? 'ol' : 'ul'}>`)
        inList = false
        listType = null
      }
      continue
    }
    
    // 水平线
    if (trimmed === '<hr>') {
      if (inList) {
        result.push(`</${listType === 'ol' ? 'ol' : 'ul'}>`)
        inList = false
        listType = null
      }
      result.push('<hr>')
      continue
    }
    
    // 标题
    if (trimmed.match(/^<h[1-6]>/)) {
      if (inList) {
        result.push(`</${listType === 'ol' ? 'ol' : 'ul'}>`)
        inList = false
        listType = null
      }
      result.push(trimmed)
      continue
    }
    
    // 任务列表
    if (trimmed.match(/^[-*+] \[[ x]\] /)) {
      if (!inList || listType !== 'task') {
        if (inList) {
          result.push(`</${listType === 'ol' ? 'ol' : 'ul'}>`)
        }
        result.push('<ul class="task-list">')
        inList = true
        listType = 'task'
      }
      const content = trimmed.replace(/^[-*+] \[([ x])\] /, '')
      const checked = trimmed.match(/\[x\]/) ? 'checked' : ''
      result.push(`<li class="task-list-item"><input type="checkbox" ${checked} disabled> ${content}</li>`)
      continue
    }
    
    // 无序列表
    if (trimmed.match(/^[-*+] /)) {
      if (!inList || listType !== 'ul') {
        if (inList) {
          result.push(`</${listType === 'ol' ? 'ol' : 'ul'}>`)
        }
        result.push('<ul>')
        inList = true
        listType = 'ul'
      }
      const content = trimmed.replace(/^[-*+] /, '')
      result.push(`<li>${content}</li>`)
      continue
    }
    
    // 有序列表
    if (trimmed.match(/^\d+\. /)) {
      if (!inList || listType !== 'ol') {
        if (inList) {
          result.push(`</${listType === 'ul' ? 'ul' : 'ul'}>`)
        }
        result.push('<ol>')
        inList = true
        listType = 'ol'
      }
      const content = trimmed.replace(/^\d+\. /, '')
      result.push(`<li>${content}</li>`)
      continue
    }
    
    // 引用
    if (trimmed.startsWith('&gt; ')) {
      if (inList) {
        result.push(`</${listType === 'ol' ? 'ol' : 'ul'}>`)
        inList = false
        listType = null
      }
      const content = trimmed.replace(/^&gt; /, '')
      result.push(`<blockquote>${content}</blockquote>`)
      continue
    }
    
    // 代码块
    if (trimmed.startsWith('<pre>')) {
      if (inList) {
        result.push(`</${listType === 'ol' ? 'ol' : 'ul'}>`)
        inList = false
        listType = null
      }
      result.push(trimmed)
      continue
    }
    
    // 普通段落
    if (inList) {
      result.push(`</${listType === 'ol' ? 'ol' : 'ul'}>`)
      inList = false
      listType = null
    }
    
    if (trimmed && !trimmed.startsWith('<')) {
      result.push(`<p>${trimmed}</p>`)
    } else {
      result.push(trimmed)
    }
  }
  
  // 关闭未闭合的列表
  if (inList) {
    result.push(`</${listType === 'ol' ? 'ol' : 'ul'}>`)
  }
  
  return result.join('\n')
}

// HTML 转义函数
function escapeHtml(text: string): string {
  const div = document.createElement('div')
  div.textContent = text
  return div.innerHTML
}

// 计算渲染后的内容
const renderedContent = computed(() => renderMarkdown(content.value))

// 更新历史记录
function updateHistory() {
  const currentContent = content.value
  // 移除当前位置之后的历史记录
  history.value = history.value.slice(0, historyIndex.value + 1)
  // 添加新记录
  history.value.push(currentContent)
  historyIndex.value = history.value.length - 1
  // 限制历史记录数量
  if (history.value.length > 50) {
    history.value.shift()
    historyIndex.value--
  }
  updateUndoRedoState()
}

function updateUndoRedoState() {
  canUndo.value = historyIndex.value > 0
  canRedo.value = historyIndex.value < history.value.length - 1
}

// 获取选中文本
function getSelection() {
  const textarea = editorRef.value
  if (!textarea) return { start: 0, end: 0, text: '' }
  
  const start = textarea.selectionStart
  const end = textarea.selectionEnd
  const text = content.value.substring(start, end)
  
  return { start, end, text }
}

// 设置选中文本
function setSelection(start: number, end: number) {
  nextTick(() => {
    const textarea = editorRef.value
    if (!textarea) return
    textarea.focus()
    textarea.setSelectionRange(start, end)
  })
}

// 插入文本
function insertText(text: string, selectInserted = false) {
  const { start, end } = getSelection()
  const before = content.value.substring(0, start)
  const after = content.value.substring(end)
  
  content.value = before + text + after
  updateHistory()
  
  if (selectInserted) {
    setSelection(start, start + text.length)
  } else {
    setSelection(start + text.length, start + text.length)
  }
}

// 包装选中文本
function wrapText(before: string, after: string) {
  const { start, end, text } = getSelection()
  const beforeText = content.value.substring(0, start)
  const afterText = content.value.substring(end)
  
  if (text) {
    // 有选中文本，包装它
    content.value = beforeText + before + text + after + afterText
    updateHistory()
    setSelection(start, end + before.length + after.length)
  } else {
    // 没有选中文本，插入标记并选中中间
    const insertText = before + after
    content.value = beforeText + insertText + afterText
    updateHistory()
    setSelection(start + before.length, start + before.length)
  }
}

// 撤销
function handleUndo() {
  if (historyIndex.value > 0) {
    historyIndex.value--
    const historyItem = history.value[historyIndex.value]
    if (historyItem !== undefined) {
      content.value = historyItem
    }
    updateUndoRedoState()
  }
}

// 重做
function handleRedo() {
  if (historyIndex.value < history.value.length - 1) {
    historyIndex.value++
    const historyItem = history.value[historyIndex.value]
    if (historyItem !== undefined) {
      content.value = historyItem
    }
    updateUndoRedoState()
  }
}

// 插入标题
function insertHeading(level: number) {
  const { start, text } = getSelection()
  const prefix = '#'.repeat(level) + ' '
  
  if (text) {
    // 如果选中了文本，在行首插入标题标记
    const lines = content.value.split('\n')
    let lineStart = 0
    let currentLine = 0
    
    for (let i = 0; i < lines.length; i++) {
      const line = lines[i]
      if (line === undefined) continue
      const lineEnd = lineStart + line.length
      if (start >= lineStart && start < lineEnd) {
        currentLine = i
        break
      }
      lineStart = lineEnd + 1
    }
    
    // 检查是否已经是标题
    const line = lines[currentLine]
    if (line !== undefined) {
      const headingMatch = line.match(/^(#{1,6})\s+/)
      if (headingMatch) {
        // 替换现有标题级别
        lines[currentLine] = prefix + line.replace(/^#{1,6}\s+/, '')
      } else {
        // 添加标题标记
        lines[currentLine] = prefix + line
      }
      
      content.value = lines.join('\n')
      updateHistory()
    }
  } else {
    // 没有选中文本，插入标题标记
    insertText(prefix)
  }
}

// 插入对齐（Markdown 本身不支持对齐，这里使用 HTML 标签）
function insertAlign(align: 'left' | 'center' | 'right' | 'justify') {
  const { start, end, text } = getSelection()
  // end 用于计算选中文本的结束位置
  const alignMap = {
    left: '<div align="left">',
    center: '<div align="center">',
    right: '<div align="right">',
    justify: '<div align="justify">',
  }
  
  if (text) {
    const wrapped = alignMap[align] + text + '</div>'
    const before = content.value.substring(0, start)
    const after = content.value.substring(end)
    content.value = before + wrapped + after
    updateHistory()
    setSelection(start, start + wrapped.length)
  } else {
    insertText(alignMap[align] + '\n\n</div>')
  }
}

// 插入列表
function insertList(type: 'unordered' | 'ordered' | 'task') {
  const { start } = getSelection()
  const lines = content.value.split('\n')
  let lineStart = 0
  let currentLine = 0
  
  for (let i = 0; i < lines.length; i++) {
    const line = lines[i]
    if (line === undefined) continue
    const lineEnd = lineStart + line.length
    if (start >= lineStart && start < lineEnd) {
      currentLine = i
      break
    }
    lineStart = lineEnd + 1
  }
  
  const prefix = type === 'unordered' ? '- ' : type === 'ordered' ? '1. ' : '- [ ] '
  const line = lines[currentLine]
  
  if (line === undefined) return
  
  // 检查是否已经是列表项
  if (line.match(/^[-*+]\s|^\d+\.\s|^[-*+]\s\[[\sx]\]\s/i)) {
    // 已经是列表，不重复添加
    return
  }
  
  lines[currentLine] = prefix + line
  content.value = lines.join('\n')
  updateHistory()
  setSelection(start + prefix.length, start + prefix.length + line.length)
}

// 插入链接
function insertLink() {
  const { start, end, text } = getSelection()
  const linkText = text || '链接文本'
  const linkUrl = 'https://'
  const markdown = `[${linkText}](${linkUrl})`
  
  if (text) {
    const before = content.value.substring(0, start)
    const after = content.value.substring(end)
    content.value = before + markdown + after
    updateHistory()
    setSelection(start + linkText.length + 3, start + markdown.length - 1)
  } else {
    insertText(markdown, true)
    // 选中 URL 部分
    nextTick(() => {
      setSelection(start + linkText.length + 3, start + markdown.length - 1)
    })
  }
}

// 插入图片
function insertImage() {
  const { start, end, text } = getSelection()
  const altText = text || '图片描述'
  const imageUrl = 'https://'
  const markdown = `![${altText}](${imageUrl})`
  
  if (text) {
    const before = content.value.substring(0, start)
    const after = content.value.substring(end)
    content.value = before + markdown + after
    updateHistory()
    setSelection(start + altText.length + 4, start + markdown.length - 1)
  } else {
    insertText(markdown, true)
    nextTick(() => {
      setSelection(start + altText.length + 4, start + markdown.length - 1)
    })
  }
}

// 插入代码块
function insertCodeBlock() {
  const { start, end, text } = getSelection()
  const codeBlock = '```\n' + (text || '代码') + '\n```'
  
  if (text) {
    const before = content.value.substring(0, start)
    const after = content.value.substring(end)
    content.value = before + codeBlock + after
    updateHistory()
    setSelection(start + 4, start + 4 + text.length)
  } else {
    insertText(codeBlock, true)
    nextTick(() => {
      setSelection(start + 4, start + 4 + 4)
    })
  }
}

// 插入水平线
function insertHorizontalRule() {
  const { start } = getSelection()
  const lines = content.value.split('\n')
  let lineStart = 0
  let currentLine = 0
  
  for (let i = 0; i < lines.length; i++) {
    const line = lines[i]
    if (line === undefined) continue
    const lineEnd = lineStart + line.length
    if (start >= lineStart && start < lineEnd) {
      currentLine = i
      break
    }
    lineStart = lineEnd + 1
  }
  
  lines.splice(currentLine + 1, 0, '---')
  content.value = lines.join('\n')
  updateHistory()
}

// 插入引用
function insertBlockquote() {
  const { start, end, text } = getSelection()
  
  if (text) {
    // 将选中的文本转换为引用
    const lines = text.split('\n')
    const quotedLines = lines.map(line => '> ' + line).join('\n')
    const before = content.value.substring(0, start)
    const after = content.value.substring(end)
    content.value = before + quotedLines + after
    updateHistory()
    setSelection(start, start + quotedLines.length)
  } else {
    insertText('> ')
  }
}

// 插入表格
function insertTable() {
  const table = `| 列1 | 列2 | 列3 |
|-----|-----|-----|
| 内容1 | 内容2 | 内容3 |
| 内容4 | 内容5 | 内容6 |`
  
  insertText(table)
}

// 处理输入
function handleInput() {
  updateHistory()
  emit('update:modelValue', content.value)
}

// 处理键盘事件
function handleKeydown(e: KeyboardEvent) {
  // Tab 键插入空格
  if (e.key === 'Tab' && !e.shiftKey) {
    e.preventDefault()
    const { start, end } = getSelection()
    const spaces = '  ' // 2个空格
    const before = content.value.substring(0, start)
    const after = content.value.substring(end)
    content.value = before + spaces + after
    updateHistory()
    setSelection(start + spaces.length, start + spaces.length)
  }
  
  // Ctrl/Cmd + Z 撤销
  if ((e.ctrlKey || e.metaKey) && e.key === 'z' && !e.shiftKey) {
    e.preventDefault()
    handleUndo()
  }
  
  // Ctrl/Cmd + Shift + Z 或 Ctrl/Cmd + Y 重做
  if ((e.ctrlKey || e.metaKey) && (e.shiftKey && e.key === 'z' || e.key === 'y')) {
    e.preventDefault()
    handleRedo()
  }
}

// 监听外部值变化
watch(
  () => props.modelValue,
  (newVal) => {
    if (newVal !== content.value) {
      content.value = newVal
      history.value = [newVal]
      historyIndex.value = 0
      updateUndoRedoState()
    }
  }
)

// 监听内容变化
watch(content, (newVal) => {
  emit('update:modelValue', newVal)
})
</script>

<style scoped lang="scss">
.markdown-editor {
  display: flex;
  flex-direction: column;
  height: 100%;
  background-color: #ffffff;

  &__toolbar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 4px 8px;
    background-color: #fafafa;
    border-bottom: 1px solid #ddd;
    // flex-wrap: wrap;
    gap: 4px;
    &-left {
        display: flex;
        align-items: center;
        gap: 4px;
    }

    &-right {
      display: flex;
      align-items: center;
      gap: 4px;
    }
  }

  &__editor-wrapper {
    flex: 1;
    min-height: 0;
    overflow: hidden;
    display: flex;
    flex-direction: column;
  }

  &__splitpanes {
    height: 100%;

    :deep(.splitpanes__pane) {
        display: flex;
        flex-direction: column;
        overflow: hidden;
    }
  }

  &__editor {
    flex: 1;
    width: 100%;
    padding: 10px;
    border: none;
    outline: none;
    resize: none;
    font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', 'Consolas', 'source-code-pro', monospace;
    font-size: 14px;
    line-height: 1.6;
    color: #333;
    background-color: #ffffff;
    overflow-y: auto;
    box-sizing: border-box;

    &::placeholder {
      color: #9ca3af;
    }
  }

  &__preview {
    flex: 1;
    padding: 10px;
    overflow-y: auto;
    background-color: #ffffff;
    color: #333;
    font-size: 14px;
    line-height: 1.8;

    &--fullscreen {
      height: 100%;
    }

    // Markdown 样式 - 使用 :deep() 穿透作用域，应用到动态插入的 HTML
    :deep(h1),
    :deep(h2),
    :deep(h3),
    :deep(h4),
    :deep(h5),
    :deep(h6) {
      margin-top: 24px;
      margin-bottom: 16px;
      font-weight: 600;
      line-height: 1.25;
    }

    :deep(h1) {
      font-size: 2em;
      border-bottom: 1px solid #eaecef;
      padding-bottom: 0.3em;
    }

    :deep(h2) {
      font-size: 1.5em;
      border-bottom: 1px solid #eaecef;
      padding-bottom: 0.3em;
    }

    :deep(h3) {
      font-size: 1.25em;
    }

    :deep(h4) {
      font-size: 1.1em;
    }

    :deep(p) {
      margin-bottom: 16px;
    }

    :deep(ul),
    :deep(ol) {
      margin-bottom: 16px;
      padding-left: 2em;
    }

    :deep(li) {
      margin-bottom: 4px;
    }

    :deep(code) {
      padding: 2px 4px;
      background-color: rgba(27, 31, 35, 0.05);
      border-radius: 3px;
      font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', 'Consolas', 'source-code-pro', monospace;
      font-size: 14px;
    }

    :deep(pre) {
      padding: 10px;
      overflow: auto;
      background-color: #f6f8fa;
      border-radius: 6px;
      margin: 8px 0;

      code {
        padding: 0;
        background-color: transparent;
        font-size: 14px;
        line-height: 1.45;
      }
    }

    :deep(blockquote) {
      padding: 0 0.5em;
      color: #6a737d;
      border-left: 0.25em solid #dfe2e5;
      margin: 0;
    }

    :deep(table) {
      border-collapse: collapse;
      margin-bottom: 16px;
      width: 100%;

      th, td {
        padding: 6px 13px;
        border: 1px solid #dfe2e5;
      }

      th {
        background-color: #f6f8fa;
        font-weight: 600;
      }
    }

    :deep(hr) {
      height: 1px;
      padding: 0;
      margin: 12px 0;
      background-color: #e1e4e8;
      border: 0;
    }

    :deep(img) {
      max-width: 100%;
      height: auto;
      margin-bottom: 16px;
    }

    :deep(a) {
      color: #0366d6;
      text-decoration: none;

      &:hover {
        text-decoration: underline;
      }
    }

    :deep(.task-list-item) {
      list-style-type: none;
      padding-left: 0;

      input[type="checkbox"] {
        margin-right: 8px;
      }
    }
  }
}

.toolbar-group {
  display: flex;
  align-items: center;
  gap: 2px;
  background-color: #eee;
  padding: 2px;
  border-radius: 4px;
}

.toolbar-divider {
  width: 1px;
  height: 24px;
  background-color: #e5e7eb;
  margin: 0 4px;
}

.toolbar-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  padding: 0;
  border: none;
  background-color: transparent;
  color: #333;
  cursor: pointer;
  border-radius: 4px;
  transition: background-color 0.2s;

  &:hover:not(:disabled) {
    background-color: #c6c6c6;
  }

  &:active:not(:disabled) {
    background-color: #cbcbcb;
  }

  &:focus {
    outline: none;
  }

  &:disabled {
    opacity: 0.4;
    cursor: not-allowed;
  }

  &.toolbar-action-btn {
    &:hover {
        color: #fff;
        background-color: #efbd9c;
    }
    &--active {
        color: #fff;
        background-color: #f2711c;
    }
  }

  &--save {
    background-color: #5a9e58;
    color: #fff;
    font-weight: 600;
    margin-right: 16px;
    
    &:hover {
      background-color: #4a8548;
      transform: scale(1.05);
    }
    
    &:active {
      background-color: #3d7340;
      transform: scale(0.95);
    }
    
    svg {
      animation: none;
    }
    
    &:hover svg {
      animation: savePulse 0.6s ease-in-out;
    }
  }

  svg {
    display: block;
  }

  strong {
    font-weight: bold;
  }

  em {
    font-style: italic;
  }
}

@keyframes savePulse {
  0%, 100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-2px);
  }
}
</style>

