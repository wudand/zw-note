<template>
  <div class="document-preview">
    <div class="document-preview__container" v-loading="pageLoading">
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
            <div class="document-preview__header">
              <h1 class="document-preview__title">
                {{ currentDocument?.title || '文档预览' }}
              </h1>
              <div class="document-preview__meta">
                <span v-if="currentDocument?.author" class="document-preview__author">
                  作者：{{ currentDocument.author }}
                </span>
                <span v-if="formattedUpdatedAt" class="document-preview__date">
                  更新时间：{{ formattedUpdatedAt }}
                </span>
              </div>
              <p
                v-if="currentDocument?.description"
                class="document-preview__desc"
              >
                {{ currentDocument.description }}
              </p>
            </div>

            <!-- 父级无正文：文档目录式子节点列表 -->
            <div v-if="viewMode === 'children'" class="document-preview__children">
              <h2 class="document-preview__section-title">{{ sectionTitle }}</h2>
              <ul class="document-preview__child-list">
                <li
                  v-for="child in childItems"
                  :key="child.id"
                  class="document-preview__child-item"
                >
                  <button
                    type="button"
                    class="document-preview__child-link"
                    @click="handleChildNavigate(child)"
                  >
                    {{ child.title }}
                  </button>
                </li>
              </ul>
            </div>

            <div
              v-else-if="!contentLoading && !content.trim()"
              class="document-preview__empty"
            >
              当前目录暂无内容
            </div>
            <div
              v-else
              class="document-preview__content"
              v-html="renderedContent"
            ></div>
          </main>
        </pane>
      </splitpanes>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount, defineAsyncComponent } from 'vue'
import { storeToRefs } from 'pinia'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import type { OutlineItem } from '@/components/DocumentOutline/index.vue'
import { useDocumentStore } from '@/store/documentStore'
import { renderMarkdown } from '@/utils/markdown'

const DocumentOutline = defineAsyncComponent(() => import('@/components/DocumentOutline/index.vue'))

const route = useRoute()
const store = useDocumentStore()
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

/** content：展示 Markdown；children：父级无正文时展示子目录入口 */
const viewMode = ref<'content' | 'children'>('content')
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
    if (outline.value.length > 0 && outline.value[0]) {
      await openOutlineItem(outline.value[0])
    }
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

  &__header {
    margin-bottom: 32px;
    padding-bottom: 24px;
    border-bottom: 1px solid #e8e8e5;
  }

  &__title {
    margin: 0 0 12px 0;
    font-size: 28px;
    font-weight: 650;
    letter-spacing: -0.03em;
    color: #111;
    line-height: 1.25;
  }

  &__meta {
    display: flex;
    gap: 20px;
    font-size: 13px;
    color: #6b6b6b;
  }

  &__desc {
    margin: 12px 0 0;
    font-size: 14px;
    line-height: 1.55;
    color: #6b6b6b;
  }

  &__author,
  &__date {
    display: flex;
    align-items: center;
  }

  &__empty {
    padding: 48px 0;
    text-align: center;
    font-size: 14px;
    color: #9b9b9b;
  }

  &__children {
    max-width: 720px;
  }

  &__section-title {
    margin: 0 0 20px;
    font-size: 22px;
    font-weight: 600;
    letter-spacing: -0.02em;
    color: #111;
    line-height: 1.3;
  }

  /* 文档 TOC 风格：圆点 + 蓝色链接，参考常见 Markdown 目录 */
  &__child-list {
    list-style: disc;
    margin: 0;
    padding: 0 0 0 1.5em;
  }

  &__child-item {
    margin: 0 0 14px;
    padding: 0;
    line-height: 1.7;
    color: #374151;

    &::marker {
      color: #4b5563;
    }

    &:last-child {
      margin-bottom: 0;
    }
  }

  &__child-link {
    display: inline;
    margin: 0;
    padding: 0;
    border: none;
    background: transparent;
    color: #3b6ea5;
    font-size: 16px;
    font-weight: 400;
    line-height: 1.7;
    letter-spacing: -0.01em;
    text-align: left;
    cursor: pointer;
    text-decoration: none;
    transition: color 160ms ease, text-decoration-color 160ms ease;

    &:hover {
      color: #2d5a8a;
      text-decoration: underline;
      text-underline-offset: 3px;
    }

    &:focus-visible {
      outline: 2px solid color-mix(in srgb, #3b6ea5 45%, transparent);
      outline-offset: 3px;
      border-radius: 2px;
    }
  }

  &__content {
    max-width: 900px;
    color: #111;
    font-size: 16px;
    line-height: 1.6;

    :deep(h1),
    :deep(h2),
    :deep(h3),
    :deep(h4),
    :deep(h5),
    :deep(h6) {
      margin-top: 32px;
      margin-bottom: 16px;
      font-weight: 600;
      line-height: 1.25;
      color: #111;
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
      line-height: 1.7;
    }

    :deep(ul),
    :deep(ol) {
      margin-bottom: 16px;
      padding-left: 2em;
    }

    :deep(li) {
      margin-bottom: 8px;
      line-height: 1.6;
    }

    :deep(code) {
      padding: 3px 6px;
      background-color: rgba(27, 31, 35, 0.05);
      border-radius: 3px;
      font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', 'Consolas', 'source-code-pro', monospace;
      font-size: 14px;
    }

    :deep(pre) {
      padding: 16px;
      overflow: auto;
      background-color: #f6f8fa;
      border-radius: 6px;
      margin: 16px 0;
      border: 1px solid #e8e8e5;

      code {
        padding: 0;
        background-color: transparent;
        font-size: 14px;
        line-height: 1.45;
      }
    }

    :deep(blockquote) {
      padding: 0 1em;
      color: #6b6b6b;
      border-left: 0.25em solid #d6d3d1;
      margin: 16px 0;
    }

    :deep(table) {
      border-collapse: collapse;
      margin-bottom: 16px;
      width: 100%;

      th,
      td {
        padding: 8px 13px;
        border: 1px solid #e8e8e5;
        word-break: break-word;
      }

      th {
        background-color: #f7f7f5;
        font-weight: 600;
      }

      tbody tr:nth-child(2n) {
        background-color: #fafafa;
      }
    }

    :deep(hr) {
      height: 1px;
      padding: 0;
      margin: 24px 0;
      background-color: #e8e8e5;
      border: 0;
    }

    :deep(img) {
      max-width: 100%;
      height: auto;
      margin: 16px 0;
    }

    :deep(a) {
      color: #5a9e58;
      text-decoration: none;

      &:hover {
        text-decoration: underline;
      }
    }

    :deep(strong) {
      font-weight: 600;
    }

    :deep(em) {
      font-style: italic;
    }
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
  .document-preview {
    &__main {
      padding: 32px 24px;
    }

    &__title {
      font-size: 24px;
    }
  }
}

@media (max-width: 768px) {
  .document-preview {
    &__main {
      padding: 24px 16px;
    }

    &__title {
      font-size: 22px;
    }

    &__meta {
      flex-direction: column;
      gap: 8px;
    }
  }
}
</style>
