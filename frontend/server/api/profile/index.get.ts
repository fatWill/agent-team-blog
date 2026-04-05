/**
 * GET /api/profile
 * 透传个人资料查询到 Go 后端
 */
export default defineEventHandler(async (event) => {
  return proxyToBackend(event, '/api/profile')
})
