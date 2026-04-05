/**
 * POST /api/pv/record
 * 透传 PV 上报到 Go 后端（无鉴权）
 */
export default defineEventHandler(async (event) => {
  const body = await readBody(event)
  return proxyToBackend(event, '/api/pv/record', {
    method: 'POST',
    body,
  })
})
