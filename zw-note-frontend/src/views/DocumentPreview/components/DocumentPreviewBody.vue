<template>
  <div class="preview-body">
    <div class="preview-body__header">
      <h1 class="preview-body__title">{{ title || '文档预览' }}</h1>
      <div class="preview-body__meta">
        <span v-if="author" class="preview-body__author">作者：{{ author }}</span>
        <span v-if="formattedUpdatedAt" class="preview-body__date">
          更新时间：{{ formattedUpdatedAt }}
        </span>
      </div>
      <p v-if="description" class="preview-body__desc">{{ description }}</p>
    </div>

    <!-- 默认未选中任何目录：展示全部目录任务列表，点击直达 -->
    <div v-if="viewMode === 'overview'" class="preview-body__overview">
      <template v-if="outline.length">
        <ul class="preview-body__overview-list">
          <li v-for="root in outline" :key="root.id" class="preview-body__overview-group">
            <button
              type="button"
              class="preview-body__overview-item preview-body__overview-item--root"
              @click="emit('child-navigate', root)"
            >
              <span class="preview-body__overview-mark" aria-hidden="true"></span>
              <span class="preview-body__overview-text" :title="root.title">{{ root.title }}</span>
            </button>

            <ul v-if="root.children?.length" class="preview-body__overview-sublist">
              <li v-for="child in root.children" :key="child.id">
                <button
                  type="button"
                  class="preview-body__overview-item preview-body__overview-item--child"
                  @click="emit('child-navigate', child)"
                >
                  <span class="preview-body__overview-mark" aria-hidden="true"></span>
                  <span class="preview-body__overview-text" :title="child.title">{{ child.title }}</span>
                </button>
              </li>
            </ul>
          </li>
        </ul>
      </template>
      <div v-else class="preview-body__empty">还没有目录</div>
    </div>

    <!-- 父级无正文：文档目录式子节点列表 -->
    <div v-else-if="viewMode === 'children'" class="preview-body__children">
      <h2 class="preview-body__section-title">{{ sectionTitle }}</h2>
      <ul class="preview-body__child-list">
        <li v-for="child in childItems" :key="child.id" class="preview-body__child-item">
          <button
            type="button"
            class="preview-body__child-link"
            @click="emit('child-navigate', child)"
          >
            {{ child.title }}
          </button>
        </li>
      </ul>
    </div>

    <div v-else-if="!contentLoading && !content.trim()" class="preview-body__empty">
      当前目录暂无内容
    </div>
    <div v-else class="preview-body__content" v-html="renderedContent"></div>
  </div>
</template>

<script setup lang="ts">
import type { OutlineItem } from '@/components/DocumentOutline/index.vue'

interface Props {
  title?: string
  author?: string
  description?: string
  formattedUpdatedAt?: string
  viewMode: 'overview' | 'content' | 'children'
  /** 全部目录树，仅 overview 视图使用 */
  outline: OutlineItem[]
  sectionTitle: string
  childItems: OutlineItem[]
  contentLoading: boolean
  content: string
  renderedContent: string
}

interface Emits {
  (e: 'child-navigate', item: OutlineItem): void
}

defineProps<Props>()
const emit = defineEmits<Emits>()
</script>

<style scoped lang="scss">
.preview-body {
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

  /* 默认目录任务列表：整行可点，弱化层级差异，强调「快速点到达」 */
  &__overview {
    max-width: 720px;
  }

  &__overview-list,
  &__overview-sublist {
    list-style: none;
    margin: 0;
    padding: 0;
  }

  &__overview-group {
    margin-bottom: 2px;
  }

  &__overview-sublist {
    padding-left: 28px;
  }

  &__overview-item {
    width: 100%;
    display: flex;
    align-items: center;
    gap: 10px;
    margin: 0;
    padding: 10px 12px;
    border: none;
    border-radius: 6px;
    background: transparent;
    color: #111;
    font-size: 15px;
    font-weight: 500;
    line-height: 1.4;
    text-align: left;
    cursor: pointer;
    transition: background-color 140ms ease, color 140ms ease;

    &:hover {
      background-color: #f2f5f0;
      color: #3a6b38;
    }

    &:focus-visible {
      outline: 2px solid color-mix(in srgb, #5a9e58 45%, transparent);
      outline-offset: -2px;
    }

    &--child {
      font-size: 14px;
      font-weight: 400;
      color: #4b5563;
    }
  }

  &__overview-mark {
    flex-shrink: 0;
    width: 6px;
    height: 6px;
    border-radius: 2px;
    background: color-mix(in srgb, #5a9e58 55%, #c8c8c3);
  }

  &__overview-item--child .preview-body__overview-mark {
    width: 5px;
    height: 5px;
    background: #c8c8c3;
  }

  &__overview-text {
    flex: 1;
    min-width: 0;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
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

@media (max-width: 768px) {
  .preview-body {
    &__title {
      font-size: 20px;
    }

    &__meta {
      flex-direction: column;
      gap: 8px;
    }

    &__header {
      margin-bottom: 20px;
      padding-bottom: 16px;
    }
  }
}
</style>
