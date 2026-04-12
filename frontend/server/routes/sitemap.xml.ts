import { defineEventHandler } from 'h3'

export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig()
  const backendUrl = config.backendUrl || 'http://127.0.0.1:8080'

  // 获取所有文章列表
  let articles: Array<{ id: number | string; updatedAt?: string; createdAt?: string }> = []
  try {
    const data = await $fetch<{ list: typeof articles }>(`${backendUrl}/api/articles?page=1&pageSize=1000`)
    articles = data?.list || []
  } catch {
    articles = []
  }

  const baseUrl = 'https://fatwill.cloud'
  const now = new Date().toISOString()

  const urls = [
    // 首页
    `  <url>
    <loc>${baseUrl}/articles</loc>
    <changefreq>daily</changefreq>
    <priority>1.0</priority>
    <lastmod>${now}</lastmod>
  </url>`,
    // 文章页
    ...articles.map(article => `  <url>
    <loc>${baseUrl}/articles/${article.id}</loc>
    <changefreq>weekly</changefreq>
    <priority>0.8</priority>
    <lastmod>${article.updatedAt || article.createdAt || now}</lastmod>
  </url>`),
  ]

  const xml = `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
${urls.join('\n')}
</urlset>`

  setHeader(event, 'Content-Type', 'application/xml; charset=utf-8')
  setHeader(event, 'Cache-Control', 'public, max-age=3600')
  return xml
})
