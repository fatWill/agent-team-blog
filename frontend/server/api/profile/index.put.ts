/**
 * PUT /api/profile
 * 透传更新个人资料到 Go 后端（需鉴权）
 */
export default defineEventHandler(async (event) => {
  const body = await readBody(event)
  return proxyToBackend(event, '/api/profile', {
    method: 'PUT',
    body,
  })
})
