/**
 * GET /api/albums/:id/photos
 * 透传相册照片列表查询到 Go 后端
 */
export default defineEventHandler(async (event) => {
  const id = event.context.params?.id
  return proxyToBackend(event, `/api/albums/${id}/photos`)
})
