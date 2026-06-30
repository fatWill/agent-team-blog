/**
 * auth-check.server 插件 — SSR 阶段确定登录态
 *
 * 直接调用 Go 后端（不经过 Nuxt server 路由），避免 Nitro 内部循环调用
 * 导致的 "statusMessage.replace is not a function" 错误。
 *
 * 后端地址读取优先级：环境变量 BACKEND_URL → 默认 http://127.0.0.1:8080
 */
export default defineNuxtPlugin(async (nuxtApp) => {
  // 仅 SSR 阶段执行；客户端 hydration 时 useState 已从 payload 中恢复
  if (import.meta.client) return

  const isLoggedIn = useState<boolean>('auth:isLoggedIn', () => false)
  const username = useState<string>('auth:username', () => '')

  try {
    // 从 runtimeConfig 读取后端地址
    const { backendUrl } = useRuntimeConfig()

    // 从当前 SSR 请求中获取 cookie（包含 auth_token）
    const event = nuxtApp.ssrContext?.event
    const cookie = event ? getRequestHeader(event, 'cookie') || '' : ''

    if (!cookie) {
      // 没有 cookie，确定是未登录
      isLoggedIn.value = false
      username.value = ''
      return
    }

    // 直接请求 Go 后端，不走 Nuxt server 路由
    const res = await $fetch<{ ok: boolean; username: string }>(
      `${backendUrl}/api/auth/check`,
      {
        headers: { cookie },
        timeout: 3000, // 3秒超时，不阻塞 SSR
      },
    )

    isLoggedIn.value = true
    username.value = res.username || ''
  } catch {
    // 任何错误（网络/401/超时）都静默处理，默认未登录
    isLoggedIn.value = false
    username.value = ''
  }
})
