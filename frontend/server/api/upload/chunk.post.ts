/**
 * POST /api/upload/chunk
 * 透传分片上传（multipart/form-data）到 Go 后端（需鉴权）
 */
export default defineEventHandler(async (event) => {
  return proxyRawToBackend(event, '/api/upload/chunk')
})
