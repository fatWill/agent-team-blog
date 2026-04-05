/**
 * GET /api/messages
 * 透传留言列表查询到 Go 后端（透传 IP header）
 */
export default defineEventHandler(async (event) => {
  const query = getQuery(event)
  const params: Record<string, string> = {}
  if (query.deviceId && typeof query.deviceId === 'string') {
    params.deviceId = query.deviceId
  }
  return proxyToBackend(event, '/api/messages', { query: params })
})
