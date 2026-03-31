import { getArticleById } from '~/server/utils/articles'

export default defineEventHandler(async (event) => {
  const id = getRouterParam(event, 'id')

  if (!id) {
    throw createError({
      statusCode: 400,
      statusMessage: '缺少文章 ID',
    })
  }

  const article = await getArticleById(id)

  if (!article) {
    throw createError({
      statusCode: 404,
      statusMessage: '文章不存在',
    })
  }

  return {
    id: article.id,
    title: article.title,
    summary: article.summary,
    content: article.content,
    createdAt: article.createdAt,
    updatedAt: article.updatedAt,
  }
})
