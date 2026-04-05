/**
 * GET /api/albums
 * 透传相册列表查询到 Go 后端
 */
export default defineEventHandler(async (event) => {
  return proxyToBackend(event, '/api/albums')
})
