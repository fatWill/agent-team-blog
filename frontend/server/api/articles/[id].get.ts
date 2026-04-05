/**
 * GET /api/articles/:id
 * 透传文章详情查询到 Go 后端
 */
export default defineEventHandler(async (event) => {
  const id = event.context.params?.id
  return proxyToBackend(event, `/api/articles/${id}`)
})
