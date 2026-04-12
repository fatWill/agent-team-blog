/**
 * Tab 路由兼容重定向中间件
 * - /home → /articles
 * - /home?tab=xxx → /xxx（如 /home?tab=tools → /tools）
 */
const tabToPath: Record<string, string> = {
  articles: '/articles',
  life: '/life',
  tools: '/tools',
  'agent-team': '/agent-team',
  guestbook: '/guestbook',
  changelog: '/changelog',
}

export default defineNuxtRouteMiddleware((to) => {
  if (to.path !== '/home') return

  // /home?tab=xxx → /xxx
  const tab = to.query.tab as string | undefined
  if (tab && tabToPath[tab]) {
    return navigateTo(tabToPath[tab], { redirectCode: 301, replace: true })
  }

  // /home → /articles
  return navigateTo('/articles', { redirectCode: 301, replace: true })
})
