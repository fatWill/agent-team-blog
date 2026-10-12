/**
 * GET /api/toys/tuyile/almanac?date=YYYY-MM-DD
 * 透传单日老黄历到 Go 后端
 */
export default defineEventHandler(event => proxyTuyileJson(event, '/api/toys/tuyile/almanac'))
