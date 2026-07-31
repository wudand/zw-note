<template>
  <div class="document-edit">
    <div class="document-edit__container">
      <splitpanes>
        <pane size="15">
          <!-- <aside class="document-edit__sidebar"> -->
            <DocumentOutline
              v-model="outline"
              :showEdit="true"
              @item-click="handleOutlineClick"
              @item-change="handleOutlineChange"
            />
          <!-- </aside>           -->
        </pane>
        <pane>
          <main class="document-edit__main">
            <div class="document-edit__editor">
              <MarkdownEditor
                v-model="content"
                placeholder="请输入文档内容..."
                height="100%"
              />
            </div>
          </main>          
        </pane>
      </splitpanes>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, defineAsyncComponent, onBeforeUnmount } from 'vue'
import { storeToRefs } from 'pinia'
import { useRoute } from 'vue-router'
import type { OutlineItem } from '@/components/DocumentOutline/index.vue'
import { useDocumentStore } from '@/store/documentStore'

const MarkdownEditor = defineAsyncComponent(() => import('@/components/MarkdownEditor/index.vue'))
const DocumentOutline = defineAsyncComponent(() => import('@/components/DocumentOutline/index.vue'))

const route = useRoute()
const store = useDocumentStore()
const { currentOutline: outline, currentContent: content } = storeToRefs(store)

function handleOutlineClick(item: OutlineItem) {
  store.fetchContent(item.id)
}

function handleOutlineChange(_items: OutlineItem[]) {
  // TODO: 对接保存目录结构接口
}

onMounted(async () => {
  const documentId = route.params.id as string
  await store.fetchOutline(documentId)
  if (outline.value.length > 0 && outline.value[0]) {
    store.fetchContent(outline.value[0].id)
  }
})

onBeforeUnmount(() => {
  store.resetCurrentDocument()
})
</script>

<style scoped lang="scss">
.document-edit {
  height: 100vh;
  background-color: #f5f5f5;
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
    background-color: #ddd;
  }
}
</style>

