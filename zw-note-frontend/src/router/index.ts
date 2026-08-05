import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'document-list',
      component: () => import('@/views/DocumentList/index.vue'),
      meta: {
        title: '文档列表'
      }
    },
    {
      path: '/documents/:id?/edit',
      name: 'document-edit',
      component: () => import('@/views/DocumentEdit/index.vue'),
      meta: {
        title: '编辑文档'
      }
    },
    {
      path: '/documents/:id?/preview',
      name: 'document-preview',
      component: () => import('@/views/DocumentPreview/index.vue'),
      meta: {
        title: '预览文档'
      }
    },
  ]
})

export default router

