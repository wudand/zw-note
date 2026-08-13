import { onMounted, onUnmounted, ref } from 'vue'

/** 移动端阅读壳断点：与路由守卫（router/index.ts）共用同一个值 */
export const MOBILE_BREAKPOINT = 768

export function isMobileViewport(breakpoint = MOBILE_BREAKPOINT) {
  return window.matchMedia(`(max-width: ${breakpoint}px)`).matches
}

/**
 * 基于 matchMedia 判断是否为移动端宽度的响应式开关。
 * 用 matchMedia 而非 UA 嗅探：窗口缩放 / 平板横竖屏 / 桌面调试都能正确响应。
 */
export function useIsMobile(breakpoint = MOBILE_BREAKPOINT) {
  const mql = window.matchMedia(`(max-width: ${breakpoint}px)`)
  const isMobile = ref(mql.matches)
  const onChange = (e: MediaQueryListEvent) => {
    isMobile.value = e.matches
  }

  onMounted(() => mql.addEventListener('change', onChange))
  onUnmounted(() => mql.removeEventListener('change', onChange))

  return isMobile
}
