/**
 * GET /api/toys/tuyile/share/:hash
 * 读取对话分享（只读）：透传到 Go 后端，保留 404 与 Cache-Control
 */
export default defineEventHandler((event) => {
  const hash = getRouterParam(event, 'hash') || ''
  if (!/^[0-9A-Za-z]{12}$/.test(hash)) {
    setResponseStatus(event, 404)
    return { error: true, statusCode: 404, statusMessage: '分享不存在或已失效' }
  }
  return proxyTuyileJson(event, `/api/toys/tuyile/share/${hash}`)
})
