/**
 * POST /api/renovation/pitfall/items/batch
 * 批量提交踩坑日记条目（需登录）
 */
export default defineEventHandler(async (event) => {
  const body = await readBody(event)
  return proxyToBackend(event, '/api/renovation/pitfall/items/batch', {
    method: 'POST',
    body,
  })
})
