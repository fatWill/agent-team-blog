/**
 * GET /api/wechat/jssdk-config
 * 透传微信 JS-SDK 签名请求到 Go 后端
 * 参数: url - 当前页面完整 URL（用于签名）
 */
export default defineEventHandler(async (event) => {
  const query = getQuery(event)
  return proxyToBackend(event, '/api/wechat/jssdk-config', {
    query: { url: String(query.url || '') },
  })
})
