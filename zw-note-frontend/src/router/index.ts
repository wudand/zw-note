import { createRouter, createWebHistory } from 'vue-router'
import { isMobileViewport } from '@/composables/useIsMobile'

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

// 手机端只做查看：即使通过分享链接 / 历史记录 / 桌面转手机深链进入编辑页，也兜底跳回预览页
router.beforeEach((to) => {
  if (to.name === 'document-edit' && isMobileViewport()) {
    return { name: 'document-preview', params: to.params }
  }
})

export default router

