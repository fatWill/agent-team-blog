/**
 * POST /api/messages
 * 透传新增留言到 Go 后端（透传 IP header）
 */
export default defineEventHandler(async (event) => {
  const body = await readBody(event)
  return proxyToBackend(event, '/api/messages', {
    method: 'POST',
    body,
  })
})
