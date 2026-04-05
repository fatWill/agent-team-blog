/**
 * DELETE /api/upload/chunk
 * 透传清理分片临时文件到 Go 后端（需鉴权）
 */
export default defineEventHandler(async (event) => {
  const body = await readBody(event)
  return proxyToBackend(event, '/api/upload/chunk', {
    method: 'DELETE',
    body,
  })
})
