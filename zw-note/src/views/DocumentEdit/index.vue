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
                :height="editorHeight"
              />
            </div>
          </main>          
        </pane>
      </splitpanes>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, defineAsyncComponent } from 'vue'
import { useRoute } from 'vue-router'
import type { OutlineItem } from '@/components/DocumentOutline/index.vue'
import { getDocumentOutline, getDocumentContent } from '@/service/api/documentList'

const MarkdownEditor = defineAsyncComponent(() => import('@/components/MarkdownEditor/index.vue'))
const DocumentOutline = defineAsyncComponent(() => import('@/components/DocumentOutline/index.vue'))
const route = useRoute()

const content = ref(``)
const outline = ref<OutlineItem[]>([])

const editorHeight = ref('100%')

function handleOutlineClick(item: OutlineItem, index: number) {
  console.log('点击目录项:', item, index)
  // TODO: 实现跳转到对应章节的逻辑
  getDocumentContentData(item.id)
}

function handleOutlineChange(items: OutlineItem[]) {
  console.log('目录变更:', items)
  // TODO: 实现保存目录结构的逻辑
}

onMounted(() => {
  const documentId = route.params.id
  console.log('加载文档:', documentId)
  // TODO: 根据文档ID加载文档内容
  getDocumentOutlineData(documentId)
})

const getDocumentOutlineData = async (documentId: string) => {
  const res = await getDocumentOutline(documentId)
  console.log('获取文档目录:', res)
  outline.value = res.data
  if (outline.value.length > 0) {
    getDocumentContentData(outline.value[0]?.id)
  }
}

const getDocumentContentData = async (outlineId: string) => {
  const res = await getDocumentContent(outlineId)
  console.log('获取文档内容:', res)
  if(res.data) {
    content.value = res.data
  }
}
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

