/**
 * GET /api/pv/trend?days=7
 * 透传 PV/UV 趋势数据到 Go 后端（需鉴权）
 */
export default defineEventHandler(async (event) => {
  const query = getQuery(event)
  return proxyToBackend(event, '/api/pv/trend', {
    query: query as Record<string, string>,
  })
})
