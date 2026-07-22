/**
 * GET /api/renovation/pitfall/items
 * 获取踩坑日记条目列表（公开）
 */
export default defineEventHandler(async (event) => {
  return proxyToBackend(event, '/api/renovation/pitfall/items')
})
