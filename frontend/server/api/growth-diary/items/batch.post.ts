/**
 * POST /api/growth-diary/items/batch
 * 批量编辑成长日记条目（需登录）
 */
export default defineEventHandler(async (event) => {
  const body = await readBody(event)
  return proxyToBackend(event, '/api/growth-diary/items/batch', {
    method: 'POST',
    body,
  })
})
