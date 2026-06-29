/**
 * GET /api/admin/wechat/access-token-status
 * 获取微信 access_token 缓存状态（需鉴权）
 */
export default defineEventHandler(async (event) => {
  return proxyToBackend(event, '/api/admin/wechat/access-token-status')
})
