/**
 * GET /api/articles
 * 透传文章列表查询到 Go 后端
 */
export default defineEventHandler(async (event) => {
  const query = getQuery(event)
  const params: Record<string, string> = {}
  if (query.title && typeof query.title === 'string') {
    params.title = query.title
  }
  return proxyToBackend(event, '/api/articles', { query: params })
})
