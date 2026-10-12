/**
 * POST /api/toys/tuyile/chat
 * AG-UI 运势聊天：SSE 流式透传到 Go 后端（不缓冲）
 */
export default defineEventHandler(event => proxyTuyileStream(event, '/api/toys/tuyile/chat'))
