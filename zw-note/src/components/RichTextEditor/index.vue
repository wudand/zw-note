<template>
  <div class="rich-text-editor">
    <Toolbar
      :editor="editorRef"
      :defaultConfig="toolbarConfig"
      :mode="mode"
      class="rich-text-editor__toolbar"
    />
    <Editor
      :defaultConfig="editorConfig"
      :mode="mode"
      v-model="valueHtml"
      @onCreated="handleCreated"
      @onChange="handleChange"
      class="rich-text-editor__editor"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, shallowRef, onBeforeUnmount, watch } from 'vue'
import '@wangeditor/editor/dist/css/style.css'
import { Editor, Toolbar } from '@wangeditor/editor-for-vue'
import type { IDomEditor, IEditorConfig, IToolbarConfig } from '@wangeditor/editor'

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

const mode = 'default'

const editorRef = shallowRef<IDomEditor>()
const valueHtml = ref(props.modelValue)

const toolbarConfig: Partial<IToolbarConfig> = {
  toolbarKeys: [
    'headerSelect',
    'bold',
    'italic',
    'underline',
    'through',
    'code',
    'sup',
    'sub',
    'clearStyle',
    '|',
    'color',
    'bgColor',
    '|',
    'fontSize',
    'fontFamily',
    'lineHeight',
    '|',
    'bulletedList',
    'numberedList',
    'todo',
    '|',
    'justifyLeft',
    'justifyRight',
    'justifyJustify',
    '|',
    'insertLink',
    'insertTable',
    'codeBlock',
    'divider',
    '|',
    'insertImage',
    'insertVideo',
    '|',
    'undo',
    'redo',
    '|',
    'fullScreen',
  ],
}

const editorConfig: Partial<IEditorConfig> = {
  placeholder: props.placeholder,
  readOnly: false,
  autoFocus: false,
  MENU_CONF: {
    uploadImage: {
      server: '/api/upload',
      fieldName: 'file',
      maxFileSize: 2 * 1024 * 1024,
      allowedFileTypes: ['image/*'],
      withCredentials: false,
      timeout: 5 * 1000,
      onBeforeUpload(file: File) {
        console.log('onBeforeUpload', file)
        return file
      },
      onProgress(progress: number) {
        console.log('onProgress', progress)
      },
      onSuccess(file: File, res: any) {
        console.log('onSuccess', file, res)
      },
      onFailed(file: File, res: any) {
        console.log('onFailed', file, res)
      },
      onError(file: File, err: any, res: any) {
        console.log('onError', file, err, res)
      },
    },
  },
}

function handleCreated(editor: IDomEditor) {
  editorRef.value = editor
}

function handleChange(editor: IDomEditor) {
  const html = editor.getHtml()
  if (html !== valueHtml.value) {
    valueHtml.value = html
  }
}

watch(
  () => props.modelValue,
  (newVal) => {
    if (newVal !== valueHtml.value) {
      valueHtml.value = newVal
    }
  }
)

watch(valueHtml, (newVal) => {
  emit('update:modelValue', newVal)
})

onBeforeUnmount(() => {
  const editor = editorRef.values
  if (editor == null) return
  editor.destroy()
})
</script>

<style scoped lang="scss">
.rich-text-editor {
  display: flex;
  flex-direction: column;
  height: 100%;
  border-left: 1px solid #e5e7eb;
  // border-radius: 6px;
  overflow: hidden;
  background-color: #ffffff;

  &__toolbar {
    flex-shrink: 0;
    border-bottom: 1px solid #e5e7eb;
    background-color: #fafafa;
  }

  &__editor {
    flex: 1;
    min-height: 0;
    overflow: hidden;
    display: flex;
    flex-direction: column;

    :deep(.w-e-text-container) {
      flex: 1;
      padding: 20px;
      overflow-y: auto;
      min-height: 0;
      background-color: #ffffff;
    }

    :deep(.w-e-text-placeholder) {
      color: #9ca3af;
    }

    :deep(.w-e-scroll) {
      height: 100%;
      display: flex;
      flex-direction: column;
    }

    :deep(.w-e-text) {
      max-width: 100%;
      min-height: 100%;
    }
  }
}
</style>

