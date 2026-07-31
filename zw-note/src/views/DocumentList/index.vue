<template>
  <div class="document-list-page">
    <div class="document-list-page__container">
      <header class="document-list-page__header">
        <h1 class="document-list-page__title">WD·笔记</h1>
        <el-button type="success" size="default" @click="handleCreate">
          创建新的文档
        </el-button>
      </header>

      <main class="document-list-page__main">
        <article
          v-for="document in documents"
          :key="document.id"
          class="document-item"
          @click="handlePreview(document)"
        >
          <div class="document-item__thumbnail">
            <div class="document-item__thumbnail-box">
              <div class="document-item__thumbnail-header">
                <div class="document-item__thumbnail-line"></div>
                <div class="document-item__thumbnail-line document-item__thumbnail-line--short"></div>
              </div>
              <div class="document-item__thumbnail-content">
                <span class="document-item__thumbnail-title">{{ document.title }}</span>
                <div class="document-item__thumbnail-meta">
                  <span class="document-item__thumbnail-author">{{ document.author }}</span>
                </div>
              </div>
              <div class="document-item__thumbnail-footer">
                <div class="document-item__thumbnail-corner"></div>
              </div>
            </div>
          </div>

          <div class="document-item__content">
            <div class="document-item__header">
              <h2 class="document-item__title">{{ document.title }}</h2>
            </div>
            <p class="document-item__description">{{ document.description }}</p>
          </div>

          <div class="document-item__actions" @click.stop>
            <el-button type="success" size="default" @click="handleEdit(document)">编辑</el-button>
            <el-dropdown @command="handleDropdownCommand">
              <el-icon class="document-item__actions-icon" size="20" color="#666">
                <MoreFilled />
              </el-icon>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item :command="{ action: 'preview', document }">查看</el-dropdown-item>
                  <el-dropdown-item :command="{ action: 'settings', document }">设置</el-dropdown-item>
                  <el-dropdown-item class="document-item__actions-delete" :command="{ action: 'delete', document }">删除</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </article>
      </main>
    </div>

    <!-- 创建新文档弹窗 -->
    <AddFile ref="addFileRef" @refresh="getDocumentListData" />
  </div>
</template>

<script setup lang="ts">
import { ref, defineAsyncComponent, onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import { useRouter } from 'vue-router'
import { MoreFilled } from '@element-plus/icons-vue'
import { ElMessageBox } from 'element-plus'
import { useDocumentStore } from '@/store/documentStore'
import type { DocumentItem } from '@/store/documentStore'

const AddFile = defineAsyncComponent(() => import('./components/AddFile.vue'))

const router = useRouter()
const addFileRef = ref()
const store = useDocumentStore()
const { documents } = storeToRefs(store)

interface DropdownPayload {
  action: 'edit' | 'preview' | 'settings' | 'delete'
  document: DocumentItem
}

onMounted(() => {
  getDocumentListData()
})

const getDocumentListData = () => store.fetchDocuments()

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
  } else if (action === 'settings') {
    addFileRef.value?.open('edit', document)
  } else if (action === 'delete') {
    ElMessageBox.confirm('确定删除该文档吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })
    .then(() => store.removeDocument(document.id))
    .catch(() => {})
  }
}
</script>

<style scoped lang="scss">
.document-list-page {
  height: 100vh;
  background: linear-gradient(to bottom, #f8f9fa 0%, #ffffff 100%);
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", "PingFang SC", "Hiragino Sans GB", "Microsoft YaHei", sans-serif;

  &__container {
    max-width: 1280px;
    margin: 0 auto;
    height: 100%;
    box-sizing: border-box;
    display: flex;
    flex-direction: column;
  }

  &__header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 16px 0;
    border-bottom: 1px solid #e9ecef;
  }

  &__title {
    margin: 0;
    font-size: 22px;
    font-weight: 600;
    color: #212529;
    letter-spacing: -0.02em;
  }

  &__main {
    display: flex;
    flex-direction: column;
    gap: 10px;
    flex: 1;
    overflow: auto;
  }
}

.document-item {
  display: flex;
  gap: 16px;
  padding: 16px;
  // background: #ffffff;
  // border: 1px solid #e9ecef;
  border-bottom: 1px solid #ddd;
  // border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;

  &:hover {
    border-color: #5a9e58;
    box-shadow: 0 2px 8px rgba(103, 194, 58, 0.1);
  }

  &:hover &__thumbnail-box {
    border-color: #5a9e58;
    box-shadow: 0 2px 8px rgba(103, 194, 58, 0.1);
  }

  &:last-child {
    margin-bottom: 0;
  }

  &__thumbnail {
    flex-shrink: 0;
  }

  &__thumbnail-box {
    width: 100px;
    height: 120px;
    background: #ffffff;
    border: 2px solid #e1e4e8;
    border-radius: 4px;
    display: flex;
    flex-direction: column;
    padding: 0;
    position: relative;
    overflow: hidden;
    transition: all 0.2s ease;

    // &:hover {
    //   border-color: #5a9e58;
    //   box-shadow: 0 2px 8px rgba(3, 102, 214, 0.1);
    // }
  }

  &__thumbnail-header {
    padding: 8px 10px 6px;
    border-bottom: 1px solid #e1e4e8;
    background: #f6f8fa;
  }

  &__thumbnail-line {
    height: 2px;
    background: #5a9e58;
    margin-bottom: 4px;
    border-radius: 1px;

    &--short {
      width: 60%;
      background: #d1d5db;
    }
  }

  &__thumbnail-content {
    flex: 1;
    padding: 10px;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
  }

  &__thumbnail-title {
    font-size: 12px;
    font-weight: 600;
    color: #24292e;
    line-height: 1.4;
    overflow: hidden;
    text-overflow: ellipsis;
    display: -webkit-box;
    -webkit-line-clamp: 3;
    line-clamp: 3;
    -webkit-box-orient: vertical;
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", "Helvetica Neue", Arial, sans-serif;
    letter-spacing: -0.01em;
  }

  &__thumbnail-meta {
    margin-top: auto;
    padding-top: 8px;
    border-top: 1px solid #f1f3f5;
  }

  &__thumbnail-author {
    font-size: 9px;
    color: #6a737d;
    font-weight: 400;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }

  &__thumbnail-footer {
    height: 8px;
    background: #f6f8fa;
    border-top: 1px solid #e1e4e8;
    position: relative;
  }

  &__thumbnail-corner {
    position: absolute;
    bottom: 0;
    right: 0;
    width: 0;
    height: 0;
    border-style: solid;
    border-width: 0 0 12px 12px;
    border-color: transparent transparent #5a9e58 transparent;
  }

  &__content {
    flex: 1;
    min-width: 0;
    display: flex;
    flex-direction: column;
    gap: 6px;
  }

  &__header {
    display: flex;
    align-items: baseline;
    gap: 12px;
    flex-wrap: wrap;
  }

  &__title {
    margin: 0;
    font-size: 17px;
    font-weight: 600;
    color: #212529;
    line-height: 1.4;
    letter-spacing: -0.01em;
  }

  &__description {
    margin: 0;
    font-size: 13px;
    line-height: 1.5;
    color: #495057;
    overflow: hidden;
    text-overflow: ellipsis;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    line-clamp: 2;
    -webkit-box-orient: vertical;
  }

  &__actions {
    flex-shrink: 0;
    transition: opacity 0.2s ease;
    &-icon {
      cursor: pointer;
      transform: rotate(90deg);
      margin-top: 6px;
      margin-left: 6px;
    }

  }
}

@media (max-width: 768px) {
  .document-list-page {
    &__container {
      padding: 16px 12px;
    }

    &__header {
      flex-direction: column;
      align-items: flex-start;
      gap: 12px;
      margin-bottom: 16px;
    }

    &__title {
      font-size: 20px;
    }
  }

  .document-item {
    gap: 12px;
    padding: 14px;

    &__thumbnail-box {
      width: 90px;
      height: 110px;
    }

    &__title {
      font-size: 16px;
    }

    &__description {
      font-size: 12px;
    }

    &__actions {
      opacity: 1;
    }
  }
}

@media (max-width: 480px) {
  .document-item {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;

    &__thumbnail {
      width: 100%;
    }

    &__thumbnail-box {
      width: 100%;
      height: 140px;
    }

    &__actions {
      width: 100%;
      display: flex;
      justify-content: flex-end;
    }
  }
}

:deep(.document-item__actions-delete) {
  color: #f56c6c;
}

:deep(.document-item__actions-delete:hover) {
  background-color: #fbdddd;
  color: #f56c6c;
  font-weight: bold;
}
</style>

