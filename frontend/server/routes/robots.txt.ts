import { defineEventHandler } from 'h3'

/**
 * GET /robots.txt — 动态生成 robots.txt
 *
 * 原为 `public/robots.txt` 静态文件，因 Sitemap 行需要站点绝对 URL，
 * 静态文件无法在运行时替换域名，故改为 Nitro 动态路由，
 * 从 runtimeConfig.public.siteUrl 读取域名。
 *
 * 抓取策略：默认拒绝，只显式 Allow 公开页面。
 * 新增公开页面时需同步在此处添加 Allow 规则 + 更新 server/routes/sitemap.xml.ts。
 */
export default defineEventHandler((event) => {
  const config = useRuntimeConfig(event)
  const siteUrl = config.public.siteUrl

  const body = `User-agent: *

# 只允许公开页面
Allow: /$
Allow: /articles
Allow: /life$
Allow: /tools$
Allow: /agent-team$
Allow: /guestbook$
Allow: /changelog$

# 屏蔽所有需要鉴权的页面
Disallow: /admin
Disallow: /login
Disallow: /api/

# 屏蔽其他非公开路由（兜底）
Disallow: /

Sitemap: ${siteUrl}/sitemap.xml
`

  setHeader(event, 'Content-Type', 'text/plain; charset=utf-8')
  setHeader(event, 'Cache-Control', 'public, max-age=3600')
  return body
})
