/**
 * POST /api/photos/:id/verify-password
 * 透传照片密码验证到 Go 后端
 */
export default defineEventHandler(async (event) => {
  const id = event.context.params?.id
  const body = await readBody(event)
  return proxyToBackend(event, `/api/photos/${id}/verify-password`, {
    method: 'POST',
    body,
  })
})
