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

  // 首次进入（SSR 或内存无状态）：请求服务端验证 cookie
  const ok = await authStore.checkAuth()

  if (!ok) {
    // 未登录，携带 redirect 参数跳转登录页
    return navigateTo(`/login?redirect=${encodeURIComponent(to.fullPath)}`, { replace: true })
  }
})
