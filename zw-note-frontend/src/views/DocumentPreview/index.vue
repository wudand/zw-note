<template>
  <div class="document-preview">
    <div class="document-preview__container">
      <splitpanes>
        <!-- 左侧：文档目录 -->
        <pane size="20" min-size="15" max-size="35">
          <aside class="document-preview__sidebar">
            <DocumentOutline :modelValue="outline" @item-click="handleOutlineClickCallback" />
          </aside>
        </pane>

        <!-- 右侧：文档内容 -->
        <pane size="80">
          <main class="document-preview__main">
            <div class="document-preview__header">
              <h1 class="document-preview__title">{{ documentTitle }}</h1>
              <div class="document-preview__meta">
                <span class="document-preview__author">作者：{{ documentAuthor }}</span>
                <span class="document-preview__date">更新时间：{{ documentDate }}</span>
              </div>
            </div>
            <div class="document-preview__content" v-html="renderedContent"></div>
          </main>
        </pane>
      </splitpanes>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, defineAsyncComponent } from 'vue'
import { useRoute } from 'vue-router'
import type { OutlineItem } from '@/components/DocumentOutline/index.vue'

const DocumentOutline = defineAsyncComponent(() => import('@/components/DocumentOutline/index.vue'))

const route = useRoute()

// 文档数据
const documentTitle = ref('Docker 容器技术快速入门')
const documentAuthor = ref('技术团队')
const documentDate = ref('2024-01-15')

const markdownContent = ref(`
> 本文档介绍 Docker 容器技术的基本使用方法

# Docker 容器技术快速入门

Docker 是一个开源的应用容器引擎，让开发者可以打包他们的应用以及依赖包到一个可移植的容器中。

## 1. 基础概念

Docker 的核心概念包括：

- **镜像（Image）**：一个只读的模板，包含了运行应用所需的代码、运行时、库等
- **容器（Container）**：镜像的运行实例
- **仓库（Registry）**：存储镜像的地方

### 1.1 镜像管理

镜像是 Docker 的核心组件之一，以下是常用的镜像操作命令。

### 1.2 容器生命周期

容器从创建到销毁的完整生命周期管理。

## 2. 常用命令

以下是 Docker 的常用命令及其说明。

### 2.1 镜像操作

#### 拉取镜像

从 Docker Hub 拉取镜像到本地：

\`\`\`bash
docker pull [镜像名称]

# 示例：拉取 nginx 镜像
docker pull nginx

# 拉取指定版本
docker pull nginx:1.21
\`\`\`

#### 查看本地镜像

\`\`\`bash
docker images

# 输出示例：
# REPOSITORY   TAG       IMAGE ID       CREATED        SIZE
# nginx        latest    605c77e624dd   2 weeks ago    141MB
\`\`\`

#### 删除镜像

\`\`\`bash
docker rmi [镜像名称或ID]

# 示例
docker rmi nginx
docker rmi 605c77e624dd
\`\`\`

### 2.2 容器操作

#### 运行容器

\`\`\`bash
docker run [选项] [镜像名称] [命令]

# 常用选项：
# -d：后台运行容器
# -p：端口映射，格式：主机端口:容器端口
# --name：指定容器名称
# -v：挂载数据卷
# -e：设置环境变量

# 示例：运行 nginx 容器
docker run -d -p 80:80 --name nginx_demo nginx
\`\`\`

#### 查看容器

\`\`\`bash
# 查看运行中的容器
docker ps

# 查看所有容器（包括已停止）
docker ps -a

# 查看容器详细信息
docker inspect [容器名称或ID]
\`\`\`

#### 启动和停止容器

\`\`\`bash
# 启动已停止的容器
docker start [容器名称或ID]

# 停止运行中的容器
docker stop [容器名称或ID]

# 重启容器
docker restart [容器名称或ID]
\`\`\`

#### 删除容器

\`\`\`bash
# 删除已停止的容器
docker rm [容器名称或ID]

# 强制删除运行中的容器
docker rm -f [容器名称或ID]

# 删除所有已停止的容器
docker container prune
\`\`\`

### 2.3 日志和调试

#### 查看容器日志

\`\`\`bash
# 查看容器日志
docker logs [容器名称或ID]

# 实时查看日志
docker logs -f [容器名称或ID]

# 查看最近的日志
docker logs --tail 100 [容器名称或ID]
\`\`\`

#### 进入容器

\`\`\`bash
# 以交互式终端进入容器
docker exec -it [容器名称或ID] /bin/bash

# 或使用 sh（如果容器中没有 bash）
docker exec -it [容器名称或ID] /bin/sh
\`\`\`

## 3. Dockerfile

Dockerfile 是用来构建镜像的文本文件，包含了一系列构建指令。

### 3.1 基本语法

\`\`\`dockerfile
# 基础镜像
FROM node:18-alpine

# 设置工作目录
WORKDIR /app

# 复制文件
COPY package*.json ./

# 运行命令
RUN npm install

# 复制源代码
COPY . .

# 暴露端口
EXPOSE 3000

# 启动命令
CMD ["npm", "start"]
\`\`\`

### 3.2 构建镜像

\`\`\`bash
# 构建镜像
docker build -t [镜像名称]:[标签] [Dockerfile所在目录]

# 示例
docker build -t myapp:latest .
\`\`\`

## 4. Docker Compose

Docker Compose 用于定义和运行多容器 Docker 应用。

### 4.1 配置文件

\`\`\`yaml
version: '3.8'

services:
  web:
    build: .
    ports:
      - "3000:3000"
    depends_on:
      - db
    environment:
      - NODE_ENV=production

  db:
    image: postgres:14
    volumes:
      - postgres_data:/var/lib/postgresql/data
    environment:
      - POSTGRES_PASSWORD=password

volumes:
  postgres_data:
\`\`\`

### 4.2 常用命令

\`\`\`bash
# 启动服务
docker-compose up -d

# 停止服务
docker-compose down

# 查看日志
docker-compose logs -f

# 重新构建
docker-compose build
\`\`\`

## 5. 最佳实践

1. **使用官方镜像**：优先使用官方维护的镜像作为基础镜像
2. **精简镜像大小**：使用 Alpine 版本，清理不必要的文件
3. **利用缓存**：合理安排 Dockerfile 指令顺序
4. **使用 .dockerignore**：排除不需要的文件
5. **安全考虑**：不在镜像中存储敏感信息

---

*文档持续更新中...*
`)

// 文档目录
const outline = ref<OutlineItem[]>([
  {
    id: '1',
    title: '1. 基础概念',
    parentId: '0',
    children: [
      { id: '1-1', title: '1.1 镜像管理', parentId: '1' },
      { id: '1-2', title: '1.2 容器生命周期', parentId: '1' },
    ]
  },
  {
    id: '2',
    title: '2. 常用命令',
    parentId: '0',
    children: [
      { id: '2-1', title: '2.1 镜像操作', parentId: '2' },
      { id: '2-2', title: '2.2 容器操作', parentId: '2' },
      { id: '2-3', title: '2.3 日志和调试', parentId: '2' },
    ]
  },
  {
    id: '3',
    title: '3. Dockerfile',
    parentId: '0',
    children: [
      { id: '3-1', title: '3.1 基本语法', parentId: '3' },
      { id: '3-2', title: '3.2 构建镜像', parentId: '3' },
    ]
  },
  {
    id: '4',
    title: '4. Docker Compose',
    parentId: '0',
    children: [
      { id: '4-1', title: '4.1 配置文件', parentId: '4' },
      { id: '4-2', title: '4.2 常用命令', parentId: '4' },
    ]
  },
  {
    id: '5',
    title: '5. 最佳实践',
    parentId: '0',
    children: [
      { id: '5-1', title: '5.1 最佳实践', parentId: '5' },
    ]
  },
])

// Markdown 渲染函数
function renderMarkdown(markdown: string): string {
  if (!markdown) return ''
  
  let html = markdown
  
  // 先处理代码块
  const codeBlocks: string[] = []
  html = html.replace(/```[\s\S]*?```/g, (match) => {
    const index = codeBlocks.length
    codeBlocks.push(match)
    return `{{MARKDOWN-CODE-BLOCK-${index}}}`
  })
  
  // 转义 HTML
  html = html
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
  
  // 标题
  html = html.replace(/^###### (.*$)/gim, '<h6>$1</h6>')
  html = html.replace(/^##### (.*$)/gim, '<h5>$1</h5>')
  html = html.replace(/^#### (.*$)/gim, '<h4>$1</h4>')
  html = html.replace(/^### (.*$)/gim, '<h3>$1</h3>')
  html = html.replace(/^## (.*$)/gim, '<h2>$1</h2>')
  html = html.replace(/^# (.*$)/gim, '<h1>$1</h1>')
  
  // 水平线
  html = html.replace(/^---$/gim, '<hr>')
  
  // 行内代码
  html = html.replace(/`([^`]+)`/g, '<code>$1</code>')
  
  // 粗体
  html = html.replace(/\*\*([^*]+)\*\*/g, '<strong>$1</strong>')
  html = html.replace(/__([^_]+)__/g, '<strong>$1</strong>')
  
  // 斜体
  html = html.replace(/(?<!\*)\*([^*]+)\*(?!\*)/g, '<em>$1</em>')
  html = html.replace(/(?<!_)_([^_]+)_(?!_)/g, '<em>$1</em>')
  
  // 删除线
  html = html.replace(/~~([^~]+)~~/g, '<del>$1</del>')
  
  // 链接
  html = html.replace(/\[([^\]]+)\]\(([^)]+)\)/g, '<a href="$2" target="_blank" rel="noopener noreferrer">$1</a>')
  
  // 图片
  html = html.replace(/!\[([^\]]*)\]\(([^)]+)\)/g, '<img src="$2" alt="$1">')
  
  // 恢复代码块
  codeBlocks.forEach((codeBlock, index) => {
    const code = codeBlock.replace(/```[\w]*\n?/g, '').replace(/```/g, '').trim()
    html = html.replace(`{{MARKDOWN-CODE-BLOCK-${index}}}`, `<pre><code>${escapeHtml(code)}</code></pre>`)
  })
  
  // 按行处理
  const lines = html.split('\n')
  const result: string[] = []
  let inList = false
  let listType: 'ul' | 'ol' | null = null
  
  for (let i = 0; i < lines.length; i++) {
    const line = lines[i]
    if (line === undefined) continue
    const trimmed = line.trim()
    
    if (!trimmed) {
      if (inList) {
        result.push(`</${listType === 'ol' ? 'ol' : 'ul'}>`)
        inList = false
        listType = null
      }
      continue
    }
    
    if (trimmed === '<hr>') {
      if (inList) {
        result.push(`</${listType === 'ol' ? 'ol' : 'ul'}>`)
        inList = false
        listType = null
      }
      result.push('<hr>')
      continue
    }
    
    if (trimmed.match(/^<h[1-6]>/)) {
      if (inList) {
        result.push(`</${listType === 'ol' ? 'ol' : 'ul'}>`)
        inList = false
        listType = null
      }
      result.push(trimmed)
      continue
    }
    
    // 无序列表
    if (trimmed.match(/^[-*+] /)) {
      if (!inList || listType !== 'ul') {
        if (inList) {
          result.push(`</${listType === 'ol' ? 'ol' : 'ul'}>`)
        }
        result.push('<ul>')
        inList = true
        listType = 'ul'
      }
      const content = trimmed.replace(/^[-*+] /, '')
      result.push(`<li>${content}</li>`)
      continue
    }
    
    // 有序列表
    if (trimmed.match(/^\d+\. /)) {
      if (!inList || listType !== 'ol') {
        if (inList) {
          result.push(`</${listType === 'ul' ? 'ul' : 'ul'}>`)
        }
        result.push('<ol>')
        inList = true
        listType = 'ol'
      }
      const content = trimmed.replace(/^\d+\. /, '')
      result.push(`<li>${content}</li>`)
      continue
    }
    
    // 引用
    if (trimmed.startsWith('&gt; ')) {
      if (inList) {
        result.push(`</${listType === 'ol' ? 'ol' : 'ul'}>`)
        inList = false
        listType = null
      }
      const content = trimmed.replace(/^&gt; /, '')
      result.push(`<blockquote>${content}</blockquote>`)
      continue
    }
    
    if (trimmed.startsWith('<pre>')) {
      if (inList) {
        result.push(`</${listType === 'ol' ? 'ol' : 'ul'}>`)
        inList = false
        listType = null
      }
      result.push(trimmed)
      continue
    }
    
    if (inList) {
      result.push(`</${listType === 'ol' ? 'ol' : 'ul'}>`)
      inList = false
      listType = null
    }
    
    if (trimmed && !trimmed.startsWith('<')) {
      result.push(`<p>${trimmed}</p>`)
    } else {
      result.push(trimmed)
    }
  }
  
  if (inList) {
    result.push(`</${listType === 'ol' ? 'ol' : 'ul'}>`)
  }
  
  return result.join('\n')
}

function escapeHtml(text: string): string {
  const div = document.createElement('div')
  div.textContent = text
  return div.innerHTML
}

const renderedContent = computed(() => renderMarkdown(markdownContent.value))

function handleOutlineClickCallback(item: OutlineItem, index: number) {
  console.log('点击目录项:', item, index)
}

onMounted(() => {
  const documentId = route.params.id
  console.log('加载文档:', documentId)
  // TODO: 根据文档ID加载文档内容
})
</script>

<style scoped lang="scss">
.document-preview {
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
    height: 100%;
    background-color: #ffffff;
    border-right: 1px solid #e1e4e8;
    overflow: hidden;
    display: flex;
    flex-direction: column;
  }

  &__main {
    flex: 1;
    min-width: 0;
    height: 100%;
    overflow-y: auto;
    background-color: #ffffff;
    padding: 40px 60px;
  }

  &__header {
    margin-bottom: 32px;
    padding-bottom: 24px;
    border-bottom: 1px solid #e1e4e8;
  }

  &__title {
    margin: 0 0 16px 0;
    font-size: 32px;
    font-weight: 700;
    color: #24292e;
    line-height: 1.25;
  }

  &__meta {
    display: flex;
    gap: 24px;
    font-size: 14px;
    color: #6a737d;
  }

  &__author,
  &__date {
    display: flex;
    align-items: center;
  }

  &__content {
    max-width: 900px;
    margin: 0 auto;
    color: #24292e;
    font-size: 16px;
    line-height: 1.6;

    // Markdown 样式
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
      color: #24292e;
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
      border: 1px solid #e1e4e8;

      code {
        padding: 0;
        background-color: transparent;
        font-size: 14px;
        line-height: 1.45;
      }
    }

    :deep(blockquote) {
      padding: 0 1em;
      color: #6a737d;
      border-left: 0.25em solid #dfe2e5;
      margin: 16px 0;
    }

    :deep(table) {
      border-collapse: collapse;
      margin-bottom: 16px;
      width: 100%;

      th, td {
        padding: 6px 13px;
        border: 1px solid #dfe2e5;
      }

      th {
        background-color: #f6f8fa;
        font-weight: 600;
      }
    }

    :deep(hr) {
      height: 1px;
      padding: 0;
      margin: 24px 0;
      background-color: #e1e4e8;
      border: 0;
    }

    :deep(img) {
      max-width: 100%;
      height: auto;
      margin: 16px 0;
    }

    :deep(a) {
      color: #0366d6;
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
    background-color: #e1e4e8;
    
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
      font-size: 28px;
    }
  }
}

@media (max-width: 768px) {
  .document-preview {
    &__main {
      padding: 24px 16px;
    }

    &__title {
      font-size: 24px;
    }

    &__meta {
      flex-direction: column;
      gap: 8px;
    }
  }
}
</style>
