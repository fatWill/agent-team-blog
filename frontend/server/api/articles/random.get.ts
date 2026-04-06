/**
 * GET /api/articles/random
 * 透传随机文章到 Go 后端
 */
export default defineEventHandler(async (event) => {
  return proxyToBackend(event, '/api/articles/random')
})
