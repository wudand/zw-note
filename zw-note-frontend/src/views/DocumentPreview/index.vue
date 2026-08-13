<template>
  <div class="document-preview">
    <!-- 桌面：左右分栏，目录常驻 -->
    <div v-if="!isMobile" class="document-preview__container" v-loading="pageLoading">
      <splitpanes>
        <pane size="20" min-size="15" max-size="35">
          <aside class="document-preview__sidebar">
            <DocumentOutline
              :modelValue="outline"
              :showEdit="false"
              :active-id="currentOutlineId"
              @item-click="handleOutlineClick"
            />
          </aside>
        </pane>

        <pane size="80">
          <main class="document-preview__main" v-loading="contentLoading">
            <DocumentPreviewBody
              :title="currentDocument?.title"
              :author="currentDocument?.author"
              :description="currentDocument?.description"
              :formatted-updated-at="formattedUpdatedAt"
              :view-mode="viewMode"
              :outline="outline"
              :section-title="sectionTitle"
              :child-items="childItems"
              :content-loading="contentLoading"
              :content="content"
              :rendered-content="renderedContent"
              @child-navigate="handleChildNavigate"
            />
          </main>
        </pane>
      </splitpanes>
    </div>

    <!-- 移动端：目录收进底部抽屉，正文全宽 -->
    <div v-else class="document-preview__mobile" v-loading="pageLoading">
      <header class="document-preview__mobile-bar">
        <button
          type="button"
          class="document-preview__mobile-back"
          aria-label="返回"
          @click="router.back()"
        >
          <el-icon :size="18"><ArrowLeft /></el-icon>
        </button>
        <span class="document-preview__mobile-title" :title="currentDocument?.title">
          {{ currentDocument?.title || '文档预览' }}
        </span>
        <button
          type="button"
          class="document-preview__mobile-toc"
          aria-label="打开目录"
          @click="drawerVisible = true"
        >
          <el-icon :size="18"><List /></el-icon>
        </button>
      </header>

      <main class="document-preview__mobile-main" v-loading="contentLoading">
        <DocumentPreviewBody
          :title="currentDocument?.title"
          :author="currentDocument?.author"
          :description="currentDocument?.description"
          :formatted-updated-at="formattedUpdatedAt"
          :view-mode="viewMode"
          :outline="outline"
          :section-title="sectionTitle"
          :child-items="childItems"
          :content-loading="contentLoading"
          :content="content"
          :rendered-content="renderedContent"
          @child-navigate="handleChildNavigate"
        />
      </main>

      <el-drawer
        v-model="drawerVisible"
        direction="ltr"
        size="78%"
        :with-header="false"
        class="document-preview__drawer"
      >
        <DocumentOutline
          :modelValue="outline"
          :showEdit="false"
          :show-header="false"
          :active-id="currentOutlineId"
          @item-click="handleOutlineClickMobile"
        />
      </el-drawer>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount, defineAsyncComponent } from 'vue'
import { storeToRefs } from 'pinia'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { ArrowLeft, List } from '@element-plus/icons-vue'
import type { OutlineItem } from '@/components/DocumentOutline/index.vue'
import { useDocumentStore } from '@/store/documentStore'
import { renderMarkdown } from '@/utils/markdown'
import { useIsMobile } from '@/composables/useIsMobile'

const DocumentOutline = defineAsyncComponent(() => import('@/components/DocumentOutline/index.vue'))
const DocumentPreviewBody = defineAsyncComponent(
  () => import('./components/DocumentPreviewBody.vue'),
)

const route = useRoute()
const router = useRouter()
const store = useDocumentStore()
const isMobile = useIsMobile()
const drawerVisible = ref(false)
const {
  currentDocument,
  currentOutline: outline,
  currentContent: content,
  currentOutlineId,
  documentLoading,
  outlineLoading,
  contentLoading,
} = storeToRefs(store)

const pageLoading = computed(() => documentLoading.value || outlineLoading.value)

/**
 * overview：默认状态，未选中任何目录，展示全部目录任务列表
 * content：展示 Markdown
 * children：父级无正文时展示子目录入口
 */
const viewMode = ref<'overview' | 'content' | 'children'>('overview')
const childItems = ref<OutlineItem[]>([])
const sectionTitle = ref('')

const formattedUpdatedAt = computed(() => {
  const iso = currentDocument.value?.updated_at
  if (!iso) return ''
  const date = new Date(iso)
  if (Number.isNaN(date.getTime())) return iso
  const y = date.getFullYear()
  const m = String(date.getMonth() + 1).padStart(2, '0')
  const d = String(date.getDate()).padStart(2, '0')
  const hh = String(date.getHours()).padStart(2, '0')
  const mm = String(date.getMinutes()).padStart(2, '0')
  return `${y}-${m}-${d} ${hh}:${mm}`
})

function findOutlineItem(id: string, items: OutlineItem[] = outline.value): OutlineItem | null {
  for (const item of items) {
    if (item.id === id) return item
    if (item.children?.length) {
      const found = findOutlineItem(id, item.children)
      if (found) return found
    }
  }
  return null
}

/**
 * 父级有正文 → 展示内容
 * 父级无正文但有子目录 → 展示子目录列表，点击跳转子节点内容
 */
async function openOutlineItem(item: OutlineItem) {
  try {
    await store.fetchContent(item.id)
    const children = item.children ?? []
    if (!content.value.trim() && children.length > 0) {
      viewMode.value = 'children'
      childItems.value = children
      sectionTitle.value = item.title
    } else {
      viewMode.value = 'content'
      childItems.value = []
      sectionTitle.value = ''
    }
  } catch (error: any) {
    ElMessage.error(error?.message || '加载内容失败')
  }
}

function handleOutlineClick(item: OutlineItem) {
  openOutlineItem(item)
}

function handleOutlineClickMobile(item: OutlineItem) {
  openOutlineItem(item)
  drawerVisible.value = false
}

function handleChildNavigate(child: OutlineItem) {
  const node = findOutlineItem(child.id) ?? child
  openOutlineItem(node)
}

const renderedContent = computed(() => renderMarkdown(content.value))

onMounted(async () => {
  const documentId = route.params.id as string
  if (!documentId) {
    ElMessage.error('缺少文档 ID')
    return
  }

  try {
    await Promise.all([
      store.fetchDocumentById(documentId),
      store.fetchOutline(documentId),
    ])
  } catch (error: any) {
    ElMessage.error(error?.message || '加载文档失败')
  }
})

onBeforeUnmount(() => {
  store.resetCurrentDocument()
})
</script>

<style scoped lang="scss">
.document-preview {
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
    height: 100%;
    overflow: hidden;
    display: flex;
    flex-direction: column;
  }

  &__main {
    flex: 1;
    min-width: 0;
    height: 100%;
    box-sizing: border-box;
    overflow-y: auto;
    background-color: #ffffff;
    padding: 40px 60px;
  }

  /* ────────── 移动端阅读壳：顶栏 + 全宽正文 + 底部抽屉目录 ────────── */
  &__mobile {
    flex: 1;
    min-height: 0;
    display: flex;
    flex-direction: column;
    background-color: #ffffff;
  }

  &__mobile-bar {
    flex-shrink: 0;
    height: 48px;
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 0 8px;
    border-bottom: 1px solid #e8e8e5;
    background-color: #ffffff;
  }

  &__mobile-back,
  &__mobile-toc {
    flex-shrink: 0;
    width: 36px;
    height: 36px;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    border: none;
    border-radius: 4px;
    background: transparent;
    color: #374151;
    cursor: pointer;
    touch-action: manipulation;

    &:active {
      background-color: rgba(0, 0, 0, 0.06);
    }
  }

  &__mobile-title {
    flex: 1;
    min-width: 0;
    text-align: center;
    font-size: 15px;
    font-weight: 600;
    color: #111;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  &__mobile-main {
    flex: 1;
    min-height: 0;
    overflow-y: auto;
    padding: 20px 16px;
  }
}

.splitpanes {
  :deep(.splitpanes__splitter) {
    background-color: #e8e8e5;

    &:hover {
      background-color: #5a9e58;
    }
  }
}

@media (max-width: 1024px) {
  .document-preview__main {
    padding: 32px 24px;
  }
}
</style>

<style lang="scss">
/*
 * el-drawer 会 teleport 到 body 之外渲染，scoped + :deep() 对 teleport 出去的
 * 节点匹配不稳定，这里用不加 scoped 的全局样式块直接兜底覆盖，避免 Element Plus
 * 默认的 .el-drawer__body padding 残留。
 */
.document-preview__drawer .el-drawer__body {
  padding: 0 !important;
  overflow: hidden;
}
</style>
