/**
 * GET /api/toys/tuyile/almanac/month?month=YYYY-MM
 * 透传月度老黄历到 Go 后端
 */
export default defineEventHandler(event => proxyTuyileJson(event, '/api/toys/tuyile/almanac/month'))
