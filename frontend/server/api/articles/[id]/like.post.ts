/**
 * POST /api/articles/:id/like
 * 透传文章点赞/取消点赞到 Go 后端
 */
export default defineEventHandler(async (event) => {
  const id = event.context.params?.id
  const body = await readBody(event)
  return proxyToBackend(event, `/api/articles/${id}/like`, {
    method: 'POST',
    body,
  })
})
