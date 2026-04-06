/**
 * POST /api/articles/:id/view
 * 透传文章阅读量记录到 Go 后端
 */
export default defineEventHandler(async (event) => {
  const id = event.context.params?.id
  const body = await readBody(event)
  return proxyToBackend(event, `/api/articles/${id}/view`, {
    method: 'POST',
    body,
  })
})
