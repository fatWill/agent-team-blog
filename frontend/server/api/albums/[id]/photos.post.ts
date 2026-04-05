/**
 * POST /api/albums/:id/photos
 * 透传上传照片（multipart/form-data）到 Go 后端（需鉴权）
 */
export default defineEventHandler(async (event) => {
  const id = event.context.params?.id
  return proxyRawToBackend(event, `/api/albums/${id}/photos`)
})
