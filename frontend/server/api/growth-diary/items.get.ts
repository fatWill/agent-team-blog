/**
 * GET /api/growth-diary/items
 * 获取成长日记条目列表（公开）
 */
export default defineEventHandler(async (event) => {
  return proxyToBackend(event, '/api/growth-diary/items')
})
