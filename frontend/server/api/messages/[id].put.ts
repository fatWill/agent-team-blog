/**
 * PUT /api/messages/:id
 * 透传修改留言到 Go 后端（透传 IP header）
 */
export default defineEventHandler(async (event) => {
  const id = event.context.params?.id
  const body = await readBody(event)
  return proxyToBackend(event, `/api/messages/${id}`, {
    method: 'PUT',
    body,
  })
})
