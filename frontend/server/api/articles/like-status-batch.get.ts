/**
 * GET /api/articles/like-status-batch
 * 透传批量点赞状态查询到 Go 后端
 */
export default defineEventHandler(async (event) => {
  const query = getQuery(event)
  const params: Record<string, string> = {}
  if (query.deviceId && typeof query.deviceId === 'string') {
    params.deviceId = query.deviceId
  }
  if (query.ids && typeof query.ids === 'string') {
    params.ids = query.ids
  }
  return proxyToBackend(event, '/api/articles/like-status-batch', { query: params })
})
