/**
 * GET /api/admin/wechat/access-token-status
 * 获取微信 access_token 缓存状态（需鉴权）
 * 字段映射：后端 hasToken/remainingSec/lastRefresh → 前端 WechatTokenStatus
 */
export default defineEventHandler(async (event) => {
  const data = await proxyToBackend(event, '/api/admin/wechat/access-token-status')
  return {
    cached: data.hasToken,
    remainingSeconds: data.remainingSec,
    lastRefreshedAt: data.lastRefresh,
  }
})
