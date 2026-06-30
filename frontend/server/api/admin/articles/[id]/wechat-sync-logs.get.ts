/**
 * GET /api/admin/articles/:id/wechat-sync-logs
 * 获取文章的微信同步日志（需鉴权）
 * 字段映射：后端 list → 前端 logs
 */
export default defineEventHandler(async (event) => {
  const id = event.context.params?.id
  const data = await proxyToBackend(event, `/api/admin/articles/${id}/wechat-sync-logs`)
  return { logs: data.list || [] }
})
