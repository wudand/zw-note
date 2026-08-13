<template>
  <div class="document-list-page">
    <div class="document-list-page__container">
      <header class="document-list-page__header">
        <div class="document-list-page__brand">
          <div class="document-list-page__mark" aria-hidden="true">
            <img src="/logo_no.png" alt="" class="document-list-page__logo" />
          </div>
          <div class="document-list-page__brand-copy">
            <div class="document-list-page__brand-row">
              <h1 class="document-list-page__title">笔记</h1>
            </div>
            <p class="document-list-page__subtitle">
              {{ documentsLoading ? '加载中…' : `${documents.length} 篇文档` }}
            </p>
          </div>
        </div>
        <el-button type="primary" @click="handleCreate">新建文档</el-button>
      </header>

      <main class="document-list-page__main" v-loading="documentsLoading">
        <div v-if="!documentsLoading && documents.length === 0" class="document-list-page__empty">
          <p class="document-list-page__empty-title">还没有文档</p>
          <p class="document-list-page__empty-desc">从一篇空白页开始整理你的想法。</p>
        </div>

        <div v-else class="document-list-page__grid">
          <article
            v-for="document in documents"
            :key="document.id"
            class="note-card"
            role="button"
            tabindex="0"
            @click="handlePreview(document)"
            @keydown.enter="handlePreview(document)"
          >
            <header class="note-card__header">
              <h2 class="note-card__title" :title="document.title">{{ document.title }}</h2>
              <div class="note-card__more-wrap" @click.stop>
                <el-dropdown
                  trigger="click"
                  @command="handleDropdownCommand"
                  @visible-change="(open) => onMenuVisibleChange(document.id, open)"
                >
                  <button
                    type="button"
                    class="note-card__more"
                    :class="{ 'is-open': menuOpenId === document.id }"
                    aria-label="更多操作"
                  >
                    <el-icon :size="16"><More /></el-icon>
                  </button>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item :command="{ action: 'edit', document }">
                        编辑
                      </el-dropdown-item>
                      <el-dropdown-item :command="{ action: 'preview', document }">
                        查看
                      </el-dropdown-item>
                      <el-dropdown-item :command="{ action: 'settings', document }">
                        设置
                      </el-dropdown-item>
                      <el-dropdown-item
                        divided
                        class="note-card__actions-delete"
                        :command="{ action: 'delete', document }"
                      >
                        删除
                      </el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
              </div>
            </header>

            <p class="note-card__desc">
              {{ document.description || '暂无描述' }}
            </p>

            <footer class="note-card__footer">
              <span v-if="document.author" class="note-card__author">{{ document.author }}</span>
              <span class="note-card__dot" v-if="document.author" aria-hidden="true">·</span>
              <time class="note-card__time" :datetime="document.updated_at">
                {{ formatRelativeTime(document.updated_at) }}
              </time>
            </footer>
          </article>
        </div>
      </main>
    </div>

    <AddFile ref="addFileRef" @refresh="getDocumentListData" />
  </div>
</template>

<script setup lang="ts">
import { ref, defineAsyncComponent, onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import { useRouter } from 'vue-router'
import { More } from '@element-plus/icons-vue'
import { ElMessageBox, ElMessage } from 'element-plus'
import { useDocumentStore } from '@/store/documentStore'
import type { DocumentItem } from '@/store/documentStore'

const AddFile = defineAsyncComponent(() => import('./components/AddFile.vue'))

const router = useRouter()
const addFileRef = ref()
const store = useDocumentStore()
const { documents, documentsLoading } = storeToRefs(store)
const menuOpenId = ref<number | string | null>(null)

interface DropdownPayload {
  action: 'edit' | 'preview' | 'settings' | 'delete'
  document: DocumentItem
}

onMounted(() => {
  getDocumentListData()
})

const getDocumentListData = () => store.fetchDocuments()

function onMenuVisibleChange(id: number | string, open: boolean) {
  menuOpenId.value = open ? id : menuOpenId.value === id ? null : menuOpenId.value
}

function formatRelativeTime(iso?: string): string {
  if (!iso) return '—'
  const date = new Date(iso)
  if (Number.isNaN(date.getTime())) return '—'

  const diffMs = Date.now() - date.getTime()
  if (diffMs < 0) return '刚刚'

  const sec = Math.floor(diffMs / 1000)
  if (sec < 60) return '刚刚'
  const min = Math.floor(sec / 60)
  if (min < 60) return `${min} 分钟前`
  const hour = Math.floor(min / 60)
  if (hour < 24) return `${hour} 小时前`
  const day = Math.floor(hour / 24)
  if (day < 30) return `${day} 天前`

  const y = date.getFullYear()
  const m = String(date.getMonth() + 1).padStart(2, '0')
  const d = String(date.getDate()).padStart(2, '0')
  return `${y}-${m}-${d}`
}

function handleCreate() {
  addFileRef.value?.open('create')
}

function handlePreview(document: DocumentItem) {
  router.push({ name: 'document-preview', params: { id: document.id } })
}

function handleEdit(document: DocumentItem) {
  router.push({ name: 'document-edit', params: { id: document.id } })
}

function handleDropdownCommand(payload: DropdownPayload) {
  const { action, document } = payload

  if (action === 'preview') {
    handlePreview(document)
  } else if (action === 'edit') {
    handleEdit(document)
  } else if (action === 'settings') {
    addFileRef.value?.open('edit', document)
  } else if (action === 'delete') {
    ElMessageBox.confirm('确定删除该文档吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })
      .then(async () => {
        try {
          await store.removeDocument(document.id)
          ElMessage.success('删除成功')
        } catch (error: any) {
          ElMessage.error(error?.message || '删除失败，请稍后重试')
        }
      })
      .catch(() => {})
  }
}
</script>

<style scoped lang="scss">
/*
  Notion calm + Linear flat
  Flat Design: 无阴影 / 1px 边框 / 排版层级驱动反馈
  产品绿仅作 CTA · focus · 标题 hover 点缀
*/
.document-list-page {
  --bg: #f7f7f5;
  --surface: #ffffff;
  --fg: #111111;
  --muted: #6b6b6b;
  --faint: #9b9b9b;
  --line: #e8e8e5;
  --line-strong: #c8c8c3;
  --accent: #5a9e58;
  --accent-soft: color-mix(in srgb, var(--accent) 18%, transparent);
  --danger: #b91c1c;
  --radius: 6px;
  --ease: 160ms ease;

  min-height: 100vh;
  height: 100vh;
  background: var(--bg);
  color: var(--fg);
  font-family: "PingFang SC", "Hiragino Sans GB", "Microsoft YaHei", ui-sans-serif, system-ui,
    sans-serif;
  -webkit-font-smoothing: antialiased;

  &__container {
    max-width: 1280px;
    width: 100%;
    margin: 0 auto;
    height: 100%;
    padding: 32px 32px 40px;
    box-sizing: border-box;
    display: flex;
    flex-direction: column;
  }

  &__header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 16px;
    margin-bottom: 20px;
    padding-bottom: 16px;
    border-bottom: 1px solid var(--line);
    flex-shrink: 0;
  }

  &__brand {
    display: flex;
    align-items: center;
    gap: 14px;
    min-width: 0;
  }

  &__mark {
    width: 40px;
    height: 40px;
    flex-shrink: 0;
    display: grid;
    place-items: center;
    border-radius: 10px;
    overflow: hidden;
  }

  &__logo {
    width: 40px;
    height: 40px;
    object-fit: contain;
    display: block;
  }

  &__brand-copy {
    min-width: 0;
    display: flex;
    flex-direction: column;
    gap: 3px;
  }

  &__brand-row {
    display: flex;
    align-items: center;
    gap: 8px;
    min-width: 0;
  }

  &__title {
    margin: 0;
    font-size: 18px;
    font-weight: 600;
    letter-spacing: -0.03em;
    line-height: 1.2;
  }

  &__subtitle {
    margin: 0;
    font-size: 12px;
    color: var(--muted);
    line-height: 1.3;
    font-variant-numeric: tabular-nums;
  }

  &__main {
    flex: 1;
    min-height: 0;
    overflow: auto;
  }

  &__grid {
    display: grid;
    grid-template-columns: repeat(3, minmax(0, 1fr));
    gap: 12px;
    padding: 4px 0 24px;
  }

  &__empty {
    min-height: 360px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 4px;
    text-align: center;
    border: 1px dashed var(--line-strong);
    border-radius: var(--radius);
    background: transparent;
    padding: 64px 24px;
  }

  &__empty-title {
    margin: 0;
    font-size: 15px;
    font-weight: 600;
    letter-spacing: -0.01em;
  }

  &__empty-desc {
    margin: 0 0 16px;
    font-size: 13px;
    color: var(--muted);
    line-height: 1.5;
  }
}

.note-card {
  display: flex;
  flex-direction: column;
  gap: 10px;
  min-height: 148px;
  padding: 16px;
  background: var(--surface);
  border: 1px solid var(--line);
  border-radius: var(--radius);
  box-shadow: none;
  cursor: pointer;
  outline: none;
  touch-action: manipulation;
  transition: border-color var(--ease), box-shadow var(--ease);

  &:hover {
    /* 保持白底，避免与页面灰底融合；只用边框 + 轻描边区分 */
    border-color: color-mix(in srgb, var(--accent) 42%, var(--line-strong));
    background: var(--surface);
    box-shadow: 0 0 0 1px color-mix(in srgb, var(--accent) 18%, transparent);
  }

  &:hover .note-card__title {
    color: var(--accent);
  }

  &:active {
    border-color: var(--accent);
    box-shadow: 0 0 0 1px var(--accent-soft);
  }

  &:focus-visible {
    border-color: var(--accent);
    background: var(--surface);
    box-shadow: 0 0 0 2px var(--accent-soft);
  }

  &:hover .note-card__more,
  &:focus-within .note-card__more {
    opacity: 1;
  }

  &__header {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 8px;
  }

  &__title {
    margin: 0;
    flex: 1;
    min-width: 0;
    font-size: 15px;
    font-weight: 600;
    line-height: 1.4;
    letter-spacing: -0.02em;
    color: var(--fg);
    overflow: hidden;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    line-clamp: 2;
    -webkit-box-orient: vertical;
    transition: color var(--ease);
  }

  &__more-wrap {
    flex-shrink: 0;
    margin: -6px -6px 0 0;
  }

  &__more {
    width: 36px;
    height: 36px;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    border: none;
    border-radius: 4px;
    background: transparent;
    color: var(--faint);
    cursor: pointer;
    opacity: 0;
    transition: opacity var(--ease), background-color var(--ease), color var(--ease);

    &:hover,
    &:focus-visible,
    &.is-open {
      opacity: 1;
      color: var(--fg);
      background: color-mix(in srgb, var(--fg) 6%, transparent);
    }

    &:focus-visible {
      box-shadow: 0 0 0 2px var(--accent-soft);
    }
  }

  &__desc {
    margin: 0;
    flex: 1;
    font-size: 13px;
    line-height: 1.55;
    color: var(--muted);
    overflow: hidden;
    display: -webkit-box;
    -webkit-line-clamp: 3;
    line-clamp: 3;
    -webkit-box-orient: vertical;
  }

  &__footer {
    display: flex;
    align-items: center;
    gap: 6px;
    margin-top: auto;
    padding-top: 4px;
    border-top: 1px solid transparent;
    min-width: 0;
  }

  &__author,
  &__time,
  &__dot {
    font-size: 12px;
    line-height: 1.4;
    color: var(--faint);
  }

  &__author {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    max-width: 46%;
  }

  &__dot {
    flex-shrink: 0;
    color: color-mix(in srgb, var(--accent) 55%, var(--faint));
  }

  &__time {
    flex-shrink: 0;
    font-variant-numeric: tabular-nums;
  }
}

@media (max-width: 1024px) {
  .document-list-page__grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 640px) {
  .document-list-page {
    &__container {
      padding: 20px 14px 28px;
    }

    &__header {
      margin-bottom: 16px;
      padding-bottom: 12px;
    }

    &__grid {
      grid-template-columns: 1fr;
      gap: 10px;
    }
  }

  .note-card__more {
    opacity: 1;
  }
}

:deep(.note-card__actions-delete) {
  color: var(--danger);
}

:deep(.note-card__actions-delete:hover) {
  background-color: #fef2f2;
  color: var(--danger);
}

@media (prefers-reduced-motion: reduce) {
  .note-card,
  .note-card__title,
  .note-card__more {
    transition: none;
  }
}
</style>
