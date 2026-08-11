<template>
  <div class="document-edit">
    <div class="document-edit__container">
      <splitpanes>
        <pane size="15">
          <DocumentOutline
            v-model="outline"
            :showEdit="true"
            :title="currentDocument?.title"
            @item-click="handleOutlineClick"
          />
        </pane>
        <pane>
          <main class="document-edit__main" v-loading="contentLoading">
            <div class="document-edit__editor">
              <MarkdownEditor
                v-model="content"
                placeholder="请输入文档内容..."
                height="100%"
                :dirty="dirty"
                :saving="contentSaving"
                :disabled="!currentOutlineId || contentLoading"
                @save="handleSave"
              />
            </div>
          </main>
        </pane>
      </splitpanes>
    </div>
  </div>
</template>

<script setup lang="ts">
import {
  ref,
  computed,
  onMounted,
  onBeforeUnmount,
  defineAsyncComponent,
} from 'vue'
import { storeToRefs } from 'pinia'
import { useRoute, onBeforeRouteLeave } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { OutlineItem } from '@/components/DocumentOutline/index.vue'
import { useDocumentStore } from '@/store/documentStore'

const MarkdownEditor = defineAsyncComponent(() => import('@/components/MarkdownEditor/index.vue'))
const DocumentOutline = defineAsyncComponent(() => import('@/components/DocumentOutline/index.vue'))

const route = useRoute()
const store = useDocumentStore()
const {
  currentDocument,
  currentOutline: outline,
  currentContent: content,
  currentOutlineId,
  contentLoading,
  contentSaving,
} = storeToRefs(store)

/** 当前节点上次成功保存（或加载）时的内容快照 */
const lastSavedContent = ref('')
const dirty = computed(() => content.value !== lastSavedContent.value)

async function loadOutlineContent(outlineId: string) {
  await store.fetchContent(outlineId)
  lastSavedContent.value = content.value
}

async function persistCurrent(): Promise<boolean> {
  if (!currentOutlineId.value) {
    ElMessage.warning('请先选择目录节点')
    return false
  }
  if (!dirty.value) {
    ElMessage.info('没有需要保存的修改')
    return true
  }

  try {
    const saved = await store.saveContent(currentOutlineId.value, content.value)
    lastSavedContent.value = saved
    ElMessage.success('保存成功')
    return true
  } catch (error: any) {
    ElMessage.error(error?.message || '保存失败，请稍后重试')
    return false
  }
}

async function handleSave() {
  await persistCurrent()
}

/**
 * 有未保存修改时：确认=保存并继续，取消=不保存继续，关闭=中止
 */
async function confirmLeaveIfDirty(message: string): Promise<'proceed' | 'abort'> {
  if (!dirty.value) return 'proceed'

  try {
    await ElMessageBox.confirm(message, '未保存的修改', {
      distinguishCancelAndClose: true,
      confirmButtonText: '保存并继续',
      cancelButtonText: '不保存',
      type: 'warning',
    })
    const ok = await persistCurrent()
    return ok ? 'proceed' : 'abort'
  } catch (action) {
    if (action === 'cancel') return 'proceed'
    return 'abort'
  }
}

async function handleOutlineClick(item: OutlineItem) {
  if (item.id === currentOutlineId.value) return

  const result = await confirmLeaveIfDirty('当前目录内容尚未保存，切换前如何处理？')
  if (result === 'abort') return

  try {
    await loadOutlineContent(item.id)
  } catch (error: any) {
    ElMessage.error(error?.message || '加载内容失败')
  }
}

function onBeforeUnload(e: BeforeUnloadEvent) {
  if (!dirty.value) return
  e.preventDefault()
  e.returnValue = ''
}

onBeforeRouteLeave(async (_to, _from, next) => {
  const result = await confirmLeaveIfDirty('当前目录内容尚未保存，离开前如何处理？')
  next(result === 'proceed')
})

onMounted(async () => {
  window.addEventListener('beforeunload', onBeforeUnload)

  const documentId = route.params.id as string
  try {
    await Promise.all([
      store.fetchDocumentById(documentId),
      store.fetchOutline(documentId),
    ])
    if (outline.value.length > 0 && outline.value[0]) {
      await loadOutlineContent(outline.value[0].id)
    }
  } catch (error: any) {
    ElMessage.error(error?.message || '加载目录失败')
  }
})

onBeforeUnmount(() => {
  window.removeEventListener('beforeunload', onBeforeUnload)
  store.resetCurrentDocument()
})
</script>

<style scoped lang="scss">
.document-edit {
  height: 100vh;
  background-color: #f7f7f5;
  overflow: hidden;
  display: flex;
  flex-direction: column;

  &__container {
    display: flex;
    flex: 1;
    min-height: 0;
  }

  &__sidebar {
    width: 280px;
    flex-shrink: 0;
    height: 100%;
    overflow: hidden;
  }

  &__main {
    flex: 1;
    min-width: 0;
    height: 100%;
    overflow: hidden;
    display: flex;
    flex-direction: column;
    background: #fff;
  }

  &__editor {
    flex: 1;
    min-height: 0;
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }
}

@media (max-width: 1024px) {
  .document-edit {
    &__container {
      flex-direction: column;
    }

    &__sidebar {
      width: 100%;
      height: 300px;
      flex-shrink: 0;
    }

    &__main {
      flex: 1;
      min-height: 0;
    }
  }
}

.splitpanes {
  :deep(.splitpanes__splitter) {
    background-color: #e8e8e5;
  }
}
</style>
