/**
 * GET /api/pv/logs?page=1&page_size=20&path=&date=
 * 透传访问日志查询到 Go 后端（需鉴权）
 */
export default defineEventHandler(async (event) => {
  const query = getQuery(event)
  return proxyToBackend(event, '/api/pv/logs', {
    query: query as Record<string, string>,
  })
})
