/**
 * POST /api/auth/logout
 * 透传登出请求到 Go 后端，透传 Set-Cookie（清除 auth_token）
 */
export default defineEventHandler(async (event) => {
  return proxyToBackend(event, '/api/auth/logout', {
    method: 'POST',
  })
})
