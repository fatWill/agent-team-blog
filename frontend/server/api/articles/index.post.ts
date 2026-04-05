/**
 * POST /api/articles
 * 透传创建文章到 Go 后端（需鉴权，cookie 自动透传）
 */
export default defineEventHandler(async (event) => {
  const body = await readBody(event)
  return proxyToBackend(event, '/api/articles', {
    method: 'POST',
    body,
  })
})
