/**
 * POST /api/theme
 * 透传主题设置到 Go 后端（透传 cookie + Set-Cookie）
 */
export default defineEventHandler(async (event) => {
  const body = await readBody(event)
  return proxyToBackend(event, '/api/theme', {
    method: 'POST',
    body,
  })
})
