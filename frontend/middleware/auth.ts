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

  // 调用服务端 check 接口验证 cookie
  const ok = await authStore.checkAuth()

  if (!ok) {
    // 未登录，携带 redirect 参数跳转登录页
    return navigateTo(`/login?redirect=${encodeURIComponent(to.fullPath)}`)
  }
})
