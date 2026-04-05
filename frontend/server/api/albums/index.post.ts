/**
 * POST /api/albums
 * 透传创建相册到 Go 后端（需鉴权）
 */
export default defineEventHandler(async (event) => {
  const body = await readBody(event)
  return proxyToBackend(event, '/api/albums', {
    method: 'POST',
    body,
  })
})
