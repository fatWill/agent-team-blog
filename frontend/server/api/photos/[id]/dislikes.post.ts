/**
 * POST /api/photos/:id/dislikes
 * 透传照片踩/取消踩到 Go 后端
 */
export default defineEventHandler(async (event) => {
  const id = event.context.params?.id
  const body = await readBody(event)
  return proxyToBackend(event, `/api/photos/${id}/dislikes`, {
    method: 'POST',
    body,
  })
})
