/**
 * DELETE /api/photos/:id
 * 透传删除照片到 Go 后端（需鉴权）
 */
export default defineEventHandler(async (event) => {
  const id = event.context.params?.id
  return proxyToBackend(event, `/api/photos/${id}`, {
    method: 'DELETE',
  })
})
