/**
 * GET /api/admin/wechat/sync-logs
 * 获取全局微信同步日志列表（需鉴权）
 * 字段映射：后端 list → 前端 logs
 */
export default defineEventHandler(async (event) => {
  const query = getQuery(event)
  const params = new URLSearchParams()
  if (query.page) params.set('page', String(query.page))
  if (query.pageSize) params.set('pageSize', String(query.pageSize))
  if (query.status) params.set('status', String(query.status))
  if (query.action) params.set('action', String(query.action))
  if (query.articleId) params.set('articleId', String(query.articleId))
  if (query.startDate) params.set('startDate', String(query.startDate))
  if (query.endDate) params.set('endDate', String(query.endDate))

  const qs = params.toString()
  const url = `/api/admin/wechat/sync-logs${qs ? `?${qs}` : ''}`
  const data = await proxyToBackend(event, url)
  return {
    logs: data.list || [],
    total: data.total || 0,
    page: data.page || 1,
    pageSize: data.pageSize || 20,
  }
})
