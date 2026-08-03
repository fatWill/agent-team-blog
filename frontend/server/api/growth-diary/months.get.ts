/**
 * GET /api/growth-diary/months
 * 获取成长日记年月列表（公开）
 */
export default defineEventHandler(async (event) => {
  return proxyToBackend(event, '/api/growth-diary/months')
})
