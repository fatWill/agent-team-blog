/**
 * POST /api/admin/articles/:id/sync-wechat
 * 手动触发同步文章到微信草稿箱（需鉴权）
 */
export default defineEventHandler(async (event) => {
  const id = event.context.params?.id
  return proxyToBackend(event, `/api/admin/articles/${id}/sync-wechat`, {
    method: 'POST',
  })
})
