/**
 * auth-check 插件 — 在页面渲染前确定登录态
 *
 * 在 SSR 阶段：通过 useRequestHeaders 转发 cookie 调用 /api/auth/check
 * 在客户端 hydration 阶段：不再重复请求（利用 useState payload 传递）
 *
 * 这确保了 SSR 输出的 HTML 中登录/未登录的 UI 状态已经正确，
 * 客户端 hydration 后不会出现闪烁。
 */
export default defineNuxtPlugin(async () => {
  const isLoggedIn = useState<boolean>('auth:isLoggedIn', () => false)
  const username = useState<string>('auth:username', () => '')

  // 客户端 hydration 时：useState 已经从 SSR payload 中恢复了正确的值
  // 不需要再发请求，直接使用 payload 中的状态即可
  if (import.meta.client) return

  // SSR 阶段：发起鉴权请求确定登录态
  try {
    const headers = useRequestHeaders(['cookie'])
    const res = await $fetch<{ ok: boolean; username: string }>('/api/auth/check', {
      headers,
    })
    isLoggedIn.value = true
    username.value = res.username || ''
  } catch {
    isLoggedIn.value = false
    username.value = ''
  }
})
