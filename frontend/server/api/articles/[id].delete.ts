import { requireAuth } from '~/server/utils/auth'
import { deleteArticle } from '~/server/utils/articles'

export default defineEventHandler(async (event) => {
  // 鉴权校验
  requireAuth(event)

  const id = getRouterParam(event, 'id')

  if (!id) {
    throw createError({
      statusCode: 400,
      statusMessage: '缺少文章 ID',
    })
  }

  const deleted = await deleteArticle(id)

  if (!deleted) {
    throw createError({
      statusCode: 404,
      statusMessage: '文章不存在',
    })
  }

  return { ok: true }
})
