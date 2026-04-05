/**
 * POST /api/upload/merge
 * 透传分片合并请求到 Go 后端（需鉴权）
 */
export default defineEventHandler(async (event) => {
  const body = await readBody(event)
  return proxyToBackend(event, '/api/upload/merge', {
    method: 'POST',
    body,
  })
})
