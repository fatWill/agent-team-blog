/**
 * GET /api/admin/articles/:id/wechat-sync-logs
 * 获取文章的微信同步日志（需鉴权）
 */
export default defineEventHandler(async (event) => {
  const id = event.context.params?.id
  return proxyToBackend(event, `/api/admin/articles/${id}/wechat-sync-logs`)
})
