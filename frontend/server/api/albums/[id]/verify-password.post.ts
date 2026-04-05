/**
 * POST /api/albums/:id/verify-password
 * 透传相册密码验证到 Go 后端
 */
export default defineEventHandler(async (event) => {
  const id = event.context.params?.id
  const body = await readBody(event)
  return proxyToBackend(event, `/api/albums/${id}/verify-password`, {
    method: 'POST',
    body,
  })
})
