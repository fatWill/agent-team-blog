/**
 * GET /api/renovation/budget/items
 * 获取预算项列表（公开）
 */
export default defineEventHandler(async (event) => {
  return proxyToBackend(event, '/api/renovation/budget/items')
})
