/**
 * POST /api/toys/tuyile/share
 * 创建对话分享快照：原样透传到 Go 后端，保留 400/413/429 状态码与 Retry-After
 */
export default defineEventHandler(event => proxyTuyileJson(event, '/api/toys/tuyile/share', 'POST'))
