import { getArticleById } from '~/server/utils/articles'

export default defineEventHandler((event) => {
  const id = getRouterParam(event, 'id')

  if (!id) {
    throw createError({
      statusCode: 400,
      statusMessage: '缺少文章 ID',
    })
  }

  const article = getArticleById(id)

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
    coverImage: article.coverImage,
    content: article.content,
    likeCount: article.likeCount,
    createdAt: article.createdAt,
    updatedAt: article.updatedAt,
  }
})
