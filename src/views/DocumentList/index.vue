<template>
  <div class="document-list-page">
    <div class="document-list-page__container">
      <header class="document-list-page__header">
        <h1 class="document-list-page__title">文档列表</h1>
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
              <span class="document-item__thumbnail-title">{{ document.title }}</span>
              <div class="document-item__thumbnail-footer">
                <span class="document-item__thumbnail-author">{{ document.author }}</span>
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
            <!-- <el-dropdown split-button type="primary" @click="handleClick" @command="handleDropdownCommand">
              编辑
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item :command="{ action: 'edit', document }">编辑</el-dropdown-item>
                  <el-dropdown-item :command="{ action: 'preview', document }">查看</el-dropdown-item>
                  <el-dropdown-item :command="{ action: 'settings', document }">设置</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown> -->

            <el-button type="primary" size="default" @click="handleEdit(document)">编辑</el-button>
          </div>
        </article>
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive } from 'vue'
import { useRouter } from 'vue-router'

interface DocumentItem {
  id: string
  title: string
  description: string
  author: string
}

const router = useRouter()
const documents = reactive<DocumentItem[]>([
  {
    id: '1',
    title: '前端笔记',
    description: '每天进步一点。日积月累。',
    author: 'DODOROWO'
  },
  {
    id: '2',
    title: '后端技术简析',
    description: '技术的瓶颈,绝不是具体的语言,框架,API接口,这些东西。',
    author: '吴烦恼'
  },
  {
    id: '2',
    title: '后端技术简析',
    description: '技术的瓶颈,绝不是具体的语言,框架,API接口,这些东西。',
    author: '吴烦恼'
  },
  {
    id: '2',
    title: '后端技术简析',
    description: '技术的瓶颈,绝不是具体的语言,框架,API接口,这些东西。',
    author: '吴烦恼'
  },
  {
    id: '2',
    title: '后端技术简析',
    description: '技术的瓶颈,绝不是具体的语言,框架,API接口,这些东西。',
    author: '吴烦恼'
  },
  {
    id: '2',
    title: '后端技术简析',
    description: '技术的瓶颈,绝不是具体的语言,框架,API接口,这些东西。',
    author: '吴烦恼'
  },
  {
    id: '2',
    title: '后端技术简析',
    description: '技术的瓶颈,绝不是具体的语言,框架,API接口,这些东西。',
    author: '吴烦恼'
  }
])

interface DropdownPayload {
  action: 'edit' | 'preview' | 'settings'
  document: DocumentItem
}

function handleCreate() {
  console.info('创建新文档')
  // TODO: 实现创建文档逻辑
}

function handlePreview(document: DocumentItem) {
  router.push({ name: 'document-preview', params: { id: document.id } })
}

function handleEdit(document: DocumentItem) {
  router.push({ name: 'document-edit', params: { id: document.id } })
}

function handleDropdownCommand(payload: DropdownPayload) {
  const { action, document } = payload

  if (action === 'edit') {
    handleEdit(document)
  } else if (action === 'preview') {
    handlePreview(document)
  } else if (action === 'settings') {
    console.info(`设置操作：${document.title}`)
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
    margin-bottom: 20px;
    padding-bottom: 16px;
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
    border-color: #67c23a;
    box-shadow: 0 2px 8px rgba(103, 194, 58, 0.1);
  }

  &:last-child {
    margin-bottom: 0;
  }

  &__thumbnail {
    flex-shrink: 0;
  }

  &__thumbnail-box {
    width: 80px;
    height: 100px;
    background: linear-gradient(135deg, #67c23a 0%, #529b2e 100%);
    border-radius: 8px;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    padding: 10px 8px;
    box-shadow: 0 2px 6px rgba(103, 194, 58, 0.25);
    position: relative;
    overflow: hidden;

    &::before {
      content: '';
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      background: linear-gradient(180deg, rgba(255, 255, 255, 0.1) 0%, transparent 100%);
      pointer-events: none;
    }
  }

  &__thumbnail-title {
    font-size: 13px;
    font-weight: 600;
    color: #ffffff;
    line-height: 1.3;
    overflow: hidden;
    text-overflow: ellipsis;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    line-clamp: 2;
    -webkit-box-orient: vertical;
    position: relative;
    z-index: 1;
  }

  &__thumbnail-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-top: auto;
    position: relative;
    z-index: 1;
  }

  &__thumbnail-author {
    font-size: 9px;
    color: rgba(255, 255, 255, 0.75);
    font-weight: 500;
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
      width: 70px;
      height: 88px;
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
      height: 120px;
    }

    &__actions {
      width: 100%;
      display: flex;
      justify-content: flex-end;
    }
  }
}
</style>

