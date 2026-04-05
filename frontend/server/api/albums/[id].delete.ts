/**
 * DELETE /api/albums/:id
 * 透传删除相册到 Go 后端（需鉴权）
 */
export default defineEventHandler(async (event) => {
  const id = event.context.params?.id
  return proxyToBackend(event, `/api/albums/${id}`, {
    method: 'DELETE',
  })
})
