import { useAuthStore } from '~/stores/auth'

/**
 * 鉴权路由守卫
 * 保护需要登录才能访问的页面（/admin 等）
 * 在服务端渲染阶段即生效，不会出现内容闪烁
 */
export default defineNuxtRouteMiddleware(async (to) => {
  // 只保护 /admin 路由
  if (!to.path.startsWith('/admin')) return

  const authStore = useAuthStore()

  // 客户端导航时：内存中已有登录状态，直接放行，无需重复请求服务端
  if (import.meta.client && authStore.isLoggedIn) return

  // 客户端导航时：如果没有 auth_token cookie，直接跳转登录页（避免无谓的网络请求）
  if (import.meta.client) {
    const hasToken = document.cookie.split(';').some(c => c.trim().startsWith('auth_token='))
    if (!hasToken) {
      return navigateTo(`/login?redirect=${encodeURIComponent(to.fullPath)}`, { replace: true })
    }
  }

  // 首次进入（SSR 或内存无状态但有 cookie）：请求服务端验证 cookie
  // SSR 阶段在 middleware 上下文里读取 headers（这里有可靠的 Nuxt 请求上下文）
  const requestHeaders = import.meta.server ? useRequestHeaders(['cookie']) : {}
  const ok = await authStore.checkAuth(requestHeaders)

  if (!ok) {
    // 未登录，携带 redirect 参数跳转登录页
    return navigateTo(`/login?redirect=${encodeURIComponent(to.fullPath)}`, { replace: true })
  }
})
