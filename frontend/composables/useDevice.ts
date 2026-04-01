/**
 * 设备检测 composable
 * 用于区分移动端/PC端，统一判断逻辑
 */
export function useDevice() {
  const isMobile = ref(false)

  function update() {
    if (import.meta.server) return
    isMobile.value = window.innerWidth < 768
  }

  onMounted(() => {
    update()
    window.addEventListener('resize', update)
  })

  onUnmounted(() => {
    window.removeEventListener('resize', update)
  })

  return { isMobile }
}

/**
 * 非响应式的移动端检测（用于命令式调用场景，如 message/Modal）
 */
export function isMobileDevice(): boolean {
  if (import.meta.server) return false
  return window.innerWidth < 768
}
