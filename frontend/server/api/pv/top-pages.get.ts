/**
 * GET /api/pv/top-pages?days=7
 * 透传 Top5 页面报表到 Go 后端（需鉴权）
 */
export default defineEventHandler(async (event) => {
  const query = getQuery(event)
  return proxyToBackend(event, '/api/pv/top-pages', {
    query: query as Record<string, string>,
  })
})
