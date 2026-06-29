/**
 * GET /api/admin/wechat/server-ip
 * 获取 VPS 出口公网 IP（需鉴权）
 */
export default defineEventHandler(async (event) => {
  return proxyToBackend(event, '/api/admin/wechat/server-ip')
})
