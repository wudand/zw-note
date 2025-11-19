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

const MarkdownEditor = defineAsyncComponent(() => import('@/components/MarkdownEditor/index.vue'))
const DocumentOutline = defineAsyncComponent(() => import('@/components/DocumentOutline/index.vue'))
const route = useRoute()

const content = ref(`
> 测试测试
**粗体** 正常 *斜体* ~~删除线~~ \`代码块\`
---
__粗体__ 正常 _斜体_ ~~删除线~~ \`代码块\`

# 请输入文档内容...

## 第二章 详细说明

### 2.1 详细说明

### 2.2 详细说明

#### 安装镜像 
\`\`\`
docker  pull [镜像名称]
例：docker pull nginx
\`\`\`
#### 查看所有容器（不加-a 只查看正在运行容器）
\`\`\`
docker ps -a
-a：查看所有容器
\`\`\`
#### 运行一个容器
\`\`\`
docker run -d -p [主机端口]:[映射端口] —name [容器名称] [镜像名称] 
-d：后台运行容器
-it: 以交互式终端运行容器

例：
docker run nginx
docker run -d -p 80:80 —name nginx_demo nginx
\`\`\`

#### 启动已停止容器
\`\`\`
docker start [容器名称] 
\`\`\`

#### 重启容器
\`\`\`
docker restart [容器名称]
\`\`\`

#### 停止容器
\`\`\`
docker stop [容器名称]
\`\`\`

#### 删除容器
\`\`\`
docker rm [容器名称] -f 
-f：强制执行
例：
docker rm nginx_demo
docker rm nginx_demo -f 
\`\`\`
`)
const outline = ref<OutlineItem[]>([
  {
    id: '1',
    title: '第一章 概述',
    parentId: '0',
    children: [
      {
        id: '1-1',
        title: '1.1 概述',
        parentId: '1',
      },
      {
        id: '1-2',
        title: '1.2 概述',
        parentId: '1',
      },
    ],
  },
  {
    id: '2',
    title: '第二章 详细说明',
    parentId: '0',
  },
])

const editorHeight = ref('100%')

function handleOutlineClick(item: OutlineItem, index: number) {
  console.log('点击目录项:', item, index)
  // TODO: 实现跳转到对应章节的逻辑
}

function handleOutlineChange(items: OutlineItem[]) {
  console.log('目录变更:', items)
  // TODO: 实现保存目录结构的逻辑
}

onMounted(() => {
  const documentId = route.params.id
  console.log('加载文档:', documentId)
  // TODO: 根据文档ID加载文档内容
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

