/**
 * PUT /api/albums/:id
 * 透传更新相册到 Go 后端（需鉴权）
 */
export default defineEventHandler(async (event) => {
  const id = event.context.params?.id
  const body = await readBody(event)
  return proxyToBackend(event, `/api/albums/${id}`, {
    method: 'PUT',
    body,
  })
})
