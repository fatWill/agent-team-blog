/**
 * GET /api/articles/:id/like-status
 * 透传文章点赞状态查询到 Go 后端
 */
export default defineEventHandler(async (event) => {
  const id = event.context.params?.id
  const query = getQuery(event)
  const params: Record<string, string> = {}
  if (query.deviceId && typeof query.deviceId === 'string') {
    params.deviceId = query.deviceId
  }
  return proxyToBackend(event, `/api/articles/${id}/like-status`, { query: params })
})
