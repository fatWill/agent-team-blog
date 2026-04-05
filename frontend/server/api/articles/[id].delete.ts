/**
 * DELETE /api/articles/:id
 * 透传删除文章到 Go 后端（需鉴权）
 */
export default defineEventHandler(async (event) => {
  const id = event.context.params?.id
  return proxyToBackend(event, `/api/articles/${id}`, {
    method: 'DELETE',
  })
})
