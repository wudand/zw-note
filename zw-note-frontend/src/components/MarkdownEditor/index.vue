<template>
  <div class="markdown-editor">
    <!-- 工具栏 -->
    <div class="markdown-editor__toolbar">
        <div class="markdown-editor__toolbar-left" @mousedown.prevent>
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
            <span
              v-if="!disabled"
              class="markdown-editor__save-status"
              :class="{
                'is-dirty': dirty && !saving,
                'is-saving': saving,
              }"
            >
              {{ saving ? '保存中…' : dirty ? '未保存' : '已保存' }}
            </span>
            <button
                type="button"
                class="toolbar-btn toolbar-btn--save"
                :class="{ 'is-dirty': dirty && !saving }"
                :disabled="disabled || saving || !dirty"
                @click="save"
                :title="saveTitle"
            >
                <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
                    <polyline points="7 10 12 15 17 10"/>
                    <line x1="12" y1="15" x2="12" y2="3"/>
                </svg>
                <span class="toolbar-btn__label">保存</span>
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
          @select="rememberSelection"
          @keyup="rememberSelection"
          @mouseup="rememberSelection"
          @click="rememberSelection"
          @paste="handlePasteImage"
          @drop="handleDropImage"
          @dragover.prevent
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
              @select="rememberSelection"
              @keyup="rememberSelection"
              @mouseup="rememberSelection"
              @click="rememberSelection"
              @paste="handlePasteImage"
              @drop="handleDropImage"
              @dragover.prevent
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

    <el-dialog
      v-model="linkDialogVisible"
      title="插入链接"
      width="420px"
      append-to-body
      :close-on-click-modal="false"
      @opened="focusLinkDialogField"
      @closed="handleLinkDialogClosed"
    >
      <el-form
        ref="linkFormRef"
        :model="linkForm"
        :rules="linkRules"
        label-position="top"
        @submit.prevent="confirmInsertLink"
      >
        <el-form-item label="显示文本" prop="text">
          <el-input
            ref="linkTextInputRef"
            v-model="linkForm.text"
            placeholder="链接显示的文字（可留空，将使用地址）"
            clearable
            maxlength="200"
            @keyup.enter="confirmInsertLink"
          />
        </el-form-item>
        <el-form-item label="链接地址" prop="url">
          <el-input
            ref="linkUrlInputRef"
            v-model="linkForm.url"
            placeholder="https://example.com"
            clearable
            @keyup.enter="confirmInsertLink"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="linkDialogVisible = false">取消</el-button>
        <el-button type="success" @click="confirmInsertLink">插入</el-button>
      </template>
    </el-dialog>

    <el-dialog
      v-model="imageDialogVisible"
      title="插入图片"
      width="440px"
      append-to-body
      :close-on-click-modal="false"
      @closed="handleImageDialogClosed"
    >
      <div class="image-insert">
        <el-radio-group v-model="imageMode" size="small">
          <el-radio-button value="upload">上传图片</el-radio-button>
          <el-radio-button value="url">图片外链</el-radio-button>
        </el-radio-group>

        <div v-if="imageMode === 'upload'" class="image-insert__upload">
          <input
            ref="imageFileInputRef"
            type="file"
            accept="image/png,image/jpeg,image/webp,image/gif"
            class="image-insert__file-input"
            @change="handleImageFileChange"
          />
          <div class="image-insert__upload-row">
            <el-button :loading="imageUploading" @click="imageFileInputRef?.click()">
              {{ imageUploadedUrl ? '重新选择' : '选择本地图片' }}
            </el-button>
            <span class="image-insert__hint">支持 jpg / png / webp / gif，最大 5MB</span>
          </div>
          <div v-if="imageUploadedUrl" class="image-insert__preview">
            <img :src="resolveAssetUrl(imageUploadedUrl)" alt="预览" />
          </div>
        </div>

        <el-input
          v-else
          v-model="imageUrlInput"
          placeholder="https://example.com/pic.png"
          clearable
        />

        <el-input
          v-model="imageAlt"
          placeholder="图片描述（可选，用于替代文本）"
          maxlength="100"
          clearable
        />
      </div>
      <template #footer>
        <el-button @click="imageDialogVisible = false">取消</el-button>
        <el-button
          type="success"
          :disabled="!canConfirmImage"
          :loading="imageUploading"
          @click="confirmInsertImage"
        >
          插入
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, watch, nextTick, computed } from 'vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules, InputInstance } from 'element-plus'
import { renderMarkdown, resolveAssetUrl } from '@/utils/markdown'
import { uploadNoteImage } from '@/service/api/documentList'

interface Props {
  modelValue?: string
  placeholder?: string
  height?: string
  /** 相对上次保存是否有未保存修改 */
  dirty?: boolean
  /** 正在请求保存接口 */
  saving?: boolean
  /** 无可用目录节点或加载中时禁用保存 */
  disabled?: boolean
}

interface Emits {
  (e: 'update:modelValue', value: string): void
  (e: 'save'): void
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: '',
  placeholder: '请输入内容...',
  height: '500px',
  dirty: false,
  saving: false,
  disabled: false,
})

const emit = defineEmits<Emits>()

const editorRef = ref<HTMLTextAreaElement>()
const content = ref(props.modelValue)
const history = ref<string[]>([props.modelValue])
const historyIndex = ref(0)
const canUndo = ref(false)
const canRedo = ref(false)
/** 点击工具栏时 textarea 会失焦，需记住上次光标位置 */
const savedSelection = ref({ start: 0, end: 0 })

// 插入链接弹窗
const linkDialogVisible = ref(false)
const linkFormRef = ref<FormInstance>()
const linkTextInputRef = ref<InputInstance>()
const linkUrlInputRef = ref<InputInstance>()
const linkForm = reactive({
  text: '',
  url: '',
})
const pendingLinkSelection = ref({ start: 0, end: 0 })

// 插入图片弹窗
const imageDialogVisible = ref(false)
const imageMode = ref<'upload' | 'url'>('upload')
const imageFileInputRef = ref<HTMLInputElement>()
const imageUploading = ref(false)
const imageUploadedUrl = ref('')
const imageUrlInput = ref('')
const imageAlt = ref('')
const pendingImageSelection = ref({ start: 0, end: 0 })
const canConfirmImage = computed(() =>
  imageMode.value === 'upload' ? !!imageUploadedUrl.value : !!imageUrlInput.value.trim(),
)

const linkRules: FormRules = {
  url: [
    { required: true, message: '请输入链接地址', trigger: 'blur' },
    {
      validator: (_rule, value: string, callback) => {
        if (!value?.trim()) {
          callback(new Error('请输入链接地址'))
          return
        }
        callback()
      },
      trigger: 'blur',
    },
  ],
}

// 编辑模式
const mode = ref<'edit' | 'preview' | 'fullscreen'>('preview')

const saveTitle = computed(() => {
  if (props.disabled) return '请先选择目录'
  if (props.saving) return '保存中…'
  if (!props.dirty) return '已保存'
  return '保存当前目录内容 (⌘S / Ctrl+S)'
})

// 切换模式
function switchMode(newMode: 'edit' | 'preview' | 'fullscreen') {
  mode.value = newMode
}

function save() {
  if (props.disabled || props.saving || !props.dirty) return
  emit('update:modelValue', content.value)
  emit('save')
}

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

// 获取选中文本（失焦时回退到上次记住的光标）
function rememberSelection() {
  const textarea = editorRef.value
  if (!textarea) return
  savedSelection.value = {
    start: textarea.selectionStart,
    end: textarea.selectionEnd,
  }
}

function getSelection() {
  const textarea = editorRef.value
  const len = content.value.length

  if (textarea && document.activeElement === textarea) {
    const start = textarea.selectionStart
    const end = textarea.selectionEnd
    savedSelection.value = { start, end }
    return { start, end, text: content.value.substring(start, end) }
  }

  const start = Math.min(Math.max(0, savedSelection.value.start), len)
  const end = Math.min(Math.max(0, savedSelection.value.end), len)
  return { start, end, text: content.value.substring(start, end) }
}

// 设置选中文本
function setSelection(start: number, end: number) {
  savedSelection.value = { start, end }
  nextTick(() => {
    const textarea = editorRef.value
    if (!textarea) return
    textarea.focus()
    textarea.setSelectionRange(start, end)
    savedSelection.value = { start, end }
  })
}

/** 根据光标位置解析当前行（含行尾光标） */
function getLineAt(pos: number) {
  const lines = content.value.split('\n')
  let lineStart = 0

  for (let i = 0; i < lines.length; i++) {
    const line = lines[i] ?? ''
    const lineEnd = lineStart + line.length
    if (pos >= lineStart && pos <= lineEnd) {
      return { lines, currentLine: i, lineStart, lineEnd, line }
    }
    lineStart = lineEnd + 1
  }

  const currentLine = Math.max(0, lines.length - 1)
  const line = lines[currentLine] ?? ''
  const lastStart = content.value.length - line.length
  return {
    lines,
    currentLine,
    lineStart: lastStart,
    lineEnd: content.value.length,
    line,
  }
}

function getListPrefix(line: string): string | null {
  const task = line.match(/^(\s*- \[[ xX]\] )/)
  if (task?.[1]) return task[1]
  const ordered = line.match(/^(\s*\d+\. )/)
  if (ordered?.[1]) return ordered[1]
  const unordered = line.match(/^(\s*[-*+] )/)
  if (unordered?.[1]) return unordered[1]
  return null
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
    const { lines, currentLine, line } = getLineAt(start)
    
    // 检查是否已经是标题
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
  const open = `<div align="${align}">`
  const close = '</div>'
  
  if (text) {
    const wrapped = open + text + close
    const before = content.value.substring(0, start)
    const after = content.value.substring(end)
    content.value = before + wrapped + after
    updateHistory()
    setSelection(start + open.length, start + open.length + text.length)
  } else {
    const snippet = `${open}\n\n${close}`
    const before = content.value.substring(0, start)
    const after = content.value.substring(end)
    content.value = before + snippet + after
    updateHistory()
    // 光标落在开标签后的空行，方便直接输入
    const cursor = start + open.length + 1
    setSelection(cursor, cursor)
  }
}

// 插入列表：在当前行写入标记，光标落在标记后方便直接输入
function insertList(type: 'unordered' | 'ordered' | 'task') {
  const { start } = getSelection()
  const { lines, currentLine, lineStart, line } = getLineAt(start)
  const prefix = type === 'unordered' ? '- ' : type === 'ordered' ? '1. ' : '- [ ] '

  const existing = getListPrefix(line)
  if (existing) {
    // 已是列表项：光标移到标记后，便于继续编辑
    setSelection(lineStart + existing.length, lineStart + line.length)
    return
  }

  if (!line.trim()) {
    lines[currentLine] = prefix
    content.value = lines.join('\n')
    updateHistory()
    setSelection(lineStart + prefix.length, lineStart + prefix.length)
    return
  }

  lines[currentLine] = prefix + line
  content.value = lines.join('\n')
  updateHistory()
  setSelection(lineStart + prefix.length, lineStart + prefix.length + line.length)
}

function isLikelyUrl(value: string): boolean {
  const v = value.trim()
  if (!v || /\s/.test(v)) return false
  if (/^(https?:\/\/|mailto:|tel:)/i.test(v)) return true
  return /^[\w.-]+\.[\w.-]+(\/[\w./?%&=+#-]*)?$/i.test(v)
}

function normalizeUrl(url: string): string {
  const trimmed = url.trim()
  if (!trimmed) return trimmed
  if (/^(https?:\/\/|mailto:|tel:|\/\/)/i.test(trimmed)) return trimmed
  if (/^[\w.-]+\.[\w.-]+/i.test(trimmed)) return `https://${trimmed}`
  return trimmed
}

function resetLinkForm() {
  linkForm.text = ''
  linkForm.url = ''
  linkFormRef.value?.clearValidate()
}

function handleLinkDialogClosed() {
  resetLinkForm()
  const { start, end } = pendingLinkSelection.value
  setSelection(start, end)
}

function focusLinkDialogField() {
  nextTick(() => {
    const target = linkForm.text.trim() ? linkUrlInputRef.value : linkTextInputRef.value
    target?.focus()
  })
}

// 插入链接：弹窗填写显示文本与地址，避免只插入空 URL 模板
function insertLink() {
  const { start, end, text } = getSelection()
  pendingLinkSelection.value = { start, end }

  const selected = text.trim()
  if (selected && isLikelyUrl(selected)) {
    linkForm.text = selected
    linkForm.url = normalizeUrl(selected)
  } else {
    linkForm.text = selected
    linkForm.url = ''
  }

  linkDialogVisible.value = true
}

async function confirmInsertLink() {
  const form = linkFormRef.value
  if (!form) return

  try {
    await form.validate()
  } catch {
    return
  }

  const url = normalizeUrl(linkForm.url)
  const text = linkForm.text.trim() || url
  if (!url) return

  const { start, end } = pendingLinkSelection.value
  const markdown = `[${text}](${url})`
  const before = content.value.substring(0, start)
  const after = content.value.substring(end)
  content.value = before + markdown + after
  updateHistory()

  const cursor = start + markdown.length
  pendingLinkSelection.value = { start: cursor, end: cursor }
  linkDialogVisible.value = false
}

// 插入图片：弹窗选择本地上传或外链，避免只插入空 URL 模板
function insertImage() {
  const { start, end, text } = getSelection()
  pendingImageSelection.value = { start, end }

  imageMode.value = 'upload'
  imageAlt.value = text.trim()
  imageUrlInput.value = ''
  imageUploadedUrl.value = ''
  imageDialogVisible.value = true
}

function handleImageDialogClosed() {
  imageMode.value = 'upload'
  imageAlt.value = ''
  imageUrlInput.value = ''
  imageUploadedUrl.value = ''
  imageUploading.value = false
  if (imageFileInputRef.value) imageFileInputRef.value.value = ''
  const { start, end } = pendingImageSelection.value
  setSelection(start, end)
}

async function handleImageFileChange(e: Event) {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return

  const allowedTypes = ['image/jpeg', 'image/png', 'image/webp', 'image/gif']
  if (!allowedTypes.includes(file.type)) {
    ElMessage.warning('仅支持 jpg / png / webp / gif 格式的图片')
    input.value = ''
    return
  }
  if (file.size > 5 * 1024 * 1024) {
    ElMessage.warning('图片大小不能超过 5MB')
    input.value = ''
    return
  }

  imageUploading.value = true
  try {
    const res: any = await uploadNoteImage(file)
    const url = res?.data?.url
    if (!url) throw new Error('上传失败')
    imageUploadedUrl.value = url
    if (!imageAlt.value.trim()) {
      imageAlt.value = file.name.replace(/\.[^.]+$/, '')
    }
  } catch (error: any) {
    ElMessage.error(error?.message || '图片上传失败，请稍后重试')
    imageUploadedUrl.value = ''
  } finally {
    imageUploading.value = false
    input.value = ''
  }
}

function confirmInsertImage() {
  const url = imageMode.value === 'upload' ? imageUploadedUrl.value : normalizeUrl(imageUrlInput.value)
  if (!url) return

  const alt = imageAlt.value.trim() || '图片描述'
  const { start, end } = pendingImageSelection.value
  const markdown = `![${alt}](${url})`
  const before = content.value.substring(0, start)
  const after = content.value.substring(end)
  content.value = before + markdown + after
  updateHistory()

  const cursor = start + markdown.length
  pendingImageSelection.value = { start: cursor, end: cursor }
  imageDialogVisible.value = false
}

/**
 * 粘贴/拖拽本地图片：先插入一段可读的上传中占位文字，上传成功后原地替换为
 * 真正的图片 Markdown；避免像浏览器默认粘贴那样把剪贴板里的文件名当纯文本
 * 插进笔记（预览区就会显示成 `!screenshot_xxx` 这种无法渲染的原始文本）。
 */
async function uploadAndInsertImage(file: File) {
  const allowedTypes = ['image/jpeg', 'image/png', 'image/webp', 'image/gif']
  if (!allowedTypes.includes(file.type)) {
    ElMessage.warning('仅支持 jpg / png / webp / gif 格式的图片')
    return
  }
  if (file.size > 5 * 1024 * 1024) {
    ElMessage.warning('图片大小不能超过 5MB')
    return
  }

  const { start, end } = getSelection()
  const before = content.value.substring(0, start)
  const after = content.value.substring(end)
  const label = (file.name || '图片').replace(/\.[^.]+$/, '') || '图片'
  const placeholder = `*⏳ 正在上传图片：${label}…*`

  content.value = before + placeholder + after
  const cursorAfterPlaceholder = start + placeholder.length
  setSelection(cursorAfterPlaceholder, cursorAfterPlaceholder)

  try {
    const res: any = await uploadNoteImage(file)
    const url = res?.data?.url
    if (!url) throw new Error('上传失败')

    const markdown = `![${label}](${url})`
    const idx = content.value.indexOf(placeholder)
    if (idx !== -1) {
      content.value = content.value.substring(0, idx) + markdown + content.value.substring(idx + placeholder.length)
      const cursor = idx + markdown.length
      setSelection(cursor, cursor)
    } else {
      content.value += markdown
    }
  } catch (error: any) {
    const idx = content.value.indexOf(placeholder)
    if (idx !== -1) {
      content.value = content.value.substring(0, idx) + content.value.substring(idx + placeholder.length)
      setSelection(idx, idx)
    }
    ElMessage.error(error?.message || '图片上传失败，请稍后重试')
  } finally {
    updateHistory()
  }
}

// 粘贴图片：拦截剪贴板中的图片数据直接上传，其余情况走浏览器默认文本粘贴
async function handlePasteImage(e: ClipboardEvent) {
  const items = e.clipboardData?.items
  if (!items || items.length === 0) return

  const imageItem = Array.from(items).find(
    (item) => item.kind === 'file' && item.type.startsWith('image/'),
  )
  if (!imageItem) return

  const file = imageItem.getAsFile()
  if (!file) return

  e.preventDefault()
  await uploadAndInsertImage(file)
}

// 拖拽图片文件到编辑区：自动上传并在光标/落点处插入
async function handleDropImage(e: DragEvent) {
  const file = Array.from(e.dataTransfer?.files ?? []).find((f) => f.type.startsWith('image/'))
  if (!file) return

  e.preventDefault()
  await uploadAndInsertImage(file)
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

// 插入水平线（在当前光标处）
function insertHorizontalRule() {
  const { start, end } = getSelection()
  const before = content.value.substring(0, start)
  const after = content.value.substring(end)

  const needsLeadingNewline = before.length > 0 && !before.endsWith('\n')
  const snippet = `${needsLeadingNewline ? '\n' : ''}\n---\n\n`
  content.value = before + snippet + after
  updateHistory()
  const cursor = before.length + snippet.length
  setSelection(cursor, cursor)
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

// 插入表格（光标处插入，并选中首个表头单元格便于直接改）
function insertTable() {
  const { start, end } = getSelection()
  const before = content.value.substring(0, start)
  const after = content.value.substring(end)
  const needsLeadingNewline = before.length > 0 && !before.endsWith('\n')
  const table = `| 列1 | 列2 | 列3 |
|-----|-----|-----|
| 内容1 | 内容2 | 内容3 |
| 内容4 | 内容5 | 内容6 |`
  const snippet = `${needsLeadingNewline ? '\n' : ''}${table}\n`
  content.value = before + snippet + after
  updateHistory()
  const cellStart = before.length + (needsLeadingNewline ? 1 : 0) + 2
  setSelection(cellStart, cellStart + '列1'.length)
}

// 处理输入
function handleInput() {
  rememberSelection()
  updateHistory()
  emit('update:modelValue', content.value)
}

/** Enter 时延续当前列表项；空列表项再按 Enter 则退出列表 */
function tryContinueListOnEnter(e: KeyboardEvent): boolean {
  const { start, end } = getSelection()
  if (start !== end) return false

  const { line, lineStart } = getLineAt(start)

  const taskMatch = line.match(/^(\s*)- \[([ xX])\] (.*)$/)
  if (taskMatch) {
    const indent = taskMatch[1] ?? ''
    const itemText = taskMatch[3] ?? ''
    e.preventDefault()
    if (!itemText.trim()) {
      // 空任务项：去掉标记，退出列表
      const before = content.value.substring(0, lineStart)
      const after = content.value.substring(lineStart + line.length)
      content.value = before + indent + after
      updateHistory()
      setSelection(lineStart + indent.length, lineStart + indent.length)
      return true
    }
    const marker = `\n${indent}- [ ] `
    const before = content.value.substring(0, start)
    const after = content.value.substring(start)
    content.value = before + marker + after
    updateHistory()
    setSelection(start + marker.length, start + marker.length)
    return true
  }

  const orderedMatch = line.match(/^(\s*)(\d+)\. (.*)$/)
  if (orderedMatch) {
    const indent = orderedMatch[1] ?? ''
    const num = parseInt(orderedMatch[2] ?? '1', 10)
    const itemText = orderedMatch[3] ?? ''
    e.preventDefault()
    if (!itemText.trim()) {
      const before = content.value.substring(0, lineStart)
      const after = content.value.substring(lineStart + line.length)
      content.value = before + indent + after
      updateHistory()
      setSelection(lineStart + indent.length, lineStart + indent.length)
      return true
    }
    const marker = `\n${indent}${num + 1}. `
    const before = content.value.substring(0, start)
    const after = content.value.substring(start)
    content.value = before + marker + after
    updateHistory()
    setSelection(start + marker.length, start + marker.length)
    return true
  }

  // 无序列表（排除任务列表）
  if (!/^(\s*)[-*+] \[/.test(line)) {
    const unorderedMatch = line.match(/^(\s*)([-*+]) (.*)$/)
    if (unorderedMatch) {
      const indent = unorderedMatch[1] ?? ''
      const bullet = unorderedMatch[2] ?? '-'
      const itemText = unorderedMatch[3] ?? ''
      e.preventDefault()
      if (!itemText.trim()) {
        const before = content.value.substring(0, lineStart)
        const after = content.value.substring(lineStart + line.length)
        content.value = before + indent + after
        updateHistory()
        setSelection(lineStart + indent.length, lineStart + indent.length)
        return true
      }
      const marker = `\n${indent}${bullet} `
      const before = content.value.substring(0, start)
      const after = content.value.substring(start)
      content.value = before + marker + after
      updateHistory()
      setSelection(start + marker.length, start + marker.length)
      return true
    }
  }

  return false
}

// 处理键盘事件
function handleKeydown(e: KeyboardEvent) {
  // Ctrl/Cmd + S 保存当前目录内容
  if ((e.ctrlKey || e.metaKey) && (e.key === 's' || e.key === 'S')) {
    e.preventDefault()
    save()
    return
  }

  // Enter 延续列表
  if (e.key === 'Enter' && !e.shiftKey && !e.ctrlKey && !e.metaKey && !e.altKey) {
    if (tryContinueListOnEnter(e)) return
  }

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
      gap: 8px;
    }
  }

  &__save-status {
    font-size: 12px;
    color: #9b9b9b;
    white-space: nowrap;
    user-select: none;

    &.is-dirty {
      color: #5a9e58;
      font-weight: 500;
    }

    &.is-saving {
      color: #6b6b6b;
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

      th,
      td {
        padding: 8px 13px;
        border: 1px solid #dfe2e5;
        word-break: break-word;
      }

      th {
        background-color: #f6f8fa;
        font-weight: 600;
      }

      tbody tr:nth-child(2n) {
        background-color: #fafbfc;
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
    width: auto;
    min-width: 32px;
    padding: 0 10px;
    gap: 6px;
    margin-right: 8px;
    background-color: #ececec;
    color: #6b6b6b;
    font-weight: 500;
    transition: background-color 160ms ease, color 160ms ease, opacity 160ms ease;

    .toolbar-btn__label {
      font-size: 13px;
      line-height: 1;
    }

    &.is-dirty {
      background-color: #5a9e58;
      color: #fff;
    }

    &:hover:not(:disabled) {
      background-color: #dcdcdc;
      color: #111;
    }

    &.is-dirty:hover:not(:disabled) {
      background-color: #4a8548;
      color: #fff;
    }

    &:disabled {
      opacity: 0.55;
      transform: none;
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

.image-insert {
  display: flex;
  flex-direction: column;
  gap: 14px;

  &__upload {
    display: flex;
    flex-direction: column;
    gap: 10px;
  }

  &__upload-row {
    display: flex;
    align-items: center;
    gap: 10px;
  }

  &__file-input {
    display: none;
  }

  &__hint {
    font-size: 12px;
    color: #9b9b9b;
  }

  &__preview {
    width: 100%;
    max-height: 180px;
    overflow: hidden;
    border: 1px solid #eee;
    border-radius: 4px;
    background: #fafafa;
    display: flex;
    align-items: center;
    justify-content: center;

    img {
      max-width: 100%;
      max-height: 180px;
      object-fit: contain;
    }
  }
}
</style>


