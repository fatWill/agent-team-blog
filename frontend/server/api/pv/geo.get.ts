/**
 * GET /api/pv/geo
 * 透传访客地理分布到 Go 后端（需管理员鉴权）
 */
export default defineEventHandler(async (event) => {
  return proxyToBackend(event, '/api/pv/geo')
})
