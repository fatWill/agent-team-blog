export default defineEventHandler(async (event) => {
  const BASE_URL = 'https://fatwill.cloud'

  // 静态页面
  const staticPages = [
    { url: `${BASE_URL}/`, priority: '1.0', changefreq: 'weekly' },
    { url: `${BASE_URL}/home`, priority: '0.9', changefreq: 'daily' },
  ]

  // 动态文章页面 - 从后端获取所有文章
  let articlePages: Array<{ url: string; priority: string; changefreq: string; lastmod?: string }> = []
  try {
    const config = useRuntimeConfig()
    const backendUrl = config.backendUrl || 'http://127.0.0.1:8080'
    const res = await $fetch<{ list: Array<{ id: number; updatedAt: string }> }>(`${backendUrl}/api/articles?page=1&pageSize=1000`)
    if (res?.list) {
      articlePages = res.list.map((article) => ({
        url: `${BASE_URL}/articles/${article.id}`,
        priority: '0.8',
        changefreq: 'monthly',
        lastmod: article.updatedAt ? new Date(article.updatedAt).toISOString().split('T')[0] : undefined,
      }))
    }
  } catch (e) {
    console.error('[sitemap] 获取文章列表失败', e)
  }

  const allPages = [...staticPages, ...articlePages]

  const urls = allPages
    .map((page) => {
      const lastmod = page.lastmod ? `\n    <lastmod>${page.lastmod}</lastmod>` : ''
      return `  <url>
    <loc>${page.url}</loc>${lastmod}
    <changefreq>${page.changefreq}</changefreq>
    <priority>${page.priority}</priority>
  </url>`
    })
    .join('\n')

  const xml = `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
${urls}
</urlset>`

  setHeader(event, 'Content-Type', 'application/xml; charset=utf-8')
  setHeader(event, 'Cache-Control', 'public, max-age=3600')
  return xml
})
