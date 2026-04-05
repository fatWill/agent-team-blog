/**
 * POST /api/photos/:id/likes
 * 透传照片点赞到 Go 后端
 */
export default defineEventHandler(async (event) => {
  const id = event.context.params?.id
  const body = await readBody(event)
  return proxyToBackend(event, `/api/photos/${id}/likes`, {
    method: 'POST',
    body,
  })
})
