/**
 * GET /api/auth/check
 * 透传鉴权检查到 Go 后端，透传 Set-Cookie（续期）
 */
export default defineEventHandler(async (event) => {
  return proxyToBackend(event, '/api/auth/check')
})
