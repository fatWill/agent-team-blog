/**
 * GET /api/pv/overview
 * 透传 PV/UV 概览数据到 Go 后端（需鉴权）
 */
export default defineEventHandler(async (event) => {
  return proxyToBackend(event, '/api/pv/overview')
})
