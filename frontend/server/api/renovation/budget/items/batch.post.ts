/**
 * POST /api/renovation/budget/items/batch
 * 批量提交预算项（需登录）
 */
export default defineEventHandler(async (event) => {
  const body = await readBody(event)
  return proxyToBackend(event, '/api/renovation/budget/items/batch', {
    method: 'POST',
    body,
  })
})
