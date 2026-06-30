/**
 * useAuth - SSR 安全的登录态管理 composable
 *
 * 解决问题：
 * - 避免 SSR hydration 闪烁（服务端就能确定登录态）
 * - 统一全站登录态读取入口
 *
 * 原理：
 * - auth_token 是 HttpOnly cookie，JS 无法直接读取
 * - 通过 callOnce + $fetch 在 SSR 阶段调用 /api/auth/check（Nuxt 自动转发 cookie）
 * - 使用 useState 跨组件共享状态，确保 SSR→客户端无闪烁
 */
export function useAuth() {
  // useState 确保 SSR 渲染的值 → 客户端 hydration 一致（不闪烁）
  const isLoggedIn = useState<boolean>('auth:isLoggedIn', () => false)
  const username = useState<string>('auth:username', () => '')

  /** 手动刷新登录态（登录/登出后调用） */
  async function refreshAuth() {
    try {
      const res = await $fetch<{ ok: boolean; username: string }>('/api/auth/check')
      isLoggedIn.value = true
      username.value = res.username || ''
    } catch {
      isLoggedIn.value = false
      username.value = ''
    }
  }

  /** 登出 */
  async function logout() {
    try {
      await $fetch('/api/auth/logout', { method: 'POST' })
    } catch { /* ignore */ }
    isLoggedIn.value = false
    username.value = ''
  }

  return {
    isLoggedIn: readonly(isLoggedIn),
    username: readonly(username),
    refreshAuth,
    logout,
  }
}
