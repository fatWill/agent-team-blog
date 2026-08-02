/**
 * GET /api/renovation/pitfall/categories
 * 获取踩坑日记分类列表（公开）
 */
export default defineEventHandler(async (event) => {
  return proxyToBackend(event, '/api/renovation/pitfall/categories')
})
