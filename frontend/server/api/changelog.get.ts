/**
 * GET /api/changelog
 * 透传更新日志查询到 Go 后端
 */
export default defineEventHandler(async (event) => {
  return proxyToBackend(event, '/api/changelog')
})
