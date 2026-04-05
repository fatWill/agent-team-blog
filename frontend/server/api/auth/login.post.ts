/**
 * POST /api/auth/login
 * 透传登录请求到 Go 后端，透传 Set-Cookie（auth_token）
 */
export default defineEventHandler(async (event) => {
  const body = await readBody(event)
  return proxyToBackend(event, '/api/auth/login', {
    method: 'POST',
    body,
  })
})
