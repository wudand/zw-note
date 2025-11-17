<template>
  <div class="document-preview">
    <h1>预览文档</h1>
    <p>文档全文预览页面 - 文档ID: {{ $route.params.id }}</p>
    <p>以全屏方式呈现富文本编辑器保存的文档内容，仅提供阅读体验</p>
  </div>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router'

const route = useRoute()
</script>

<style scoped>
.document-preview {
  padding: 20px;
}
</style>

