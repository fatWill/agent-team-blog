/**
 * PUT /api/photos/:id
 * 透传更新照片信息到 Go 后端（需鉴权）
 */
export default defineEventHandler(async (event) => {
  const id = event.context.params?.id
  const body = await readBody(event)
  return proxyToBackend(event, `/api/photos/${id}`, {
    method: 'PUT',
    body,
  })
})
