/**
 * GET /api/theme
 * 透传主题查询到 Go 后端（透传 cookie 含 blog_uid）
 */
export default defineEventHandler(async (event) => {
  return proxyToBackend(event, '/api/theme')
})
