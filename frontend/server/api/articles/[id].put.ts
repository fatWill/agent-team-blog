import { readBody } from 'h3'
import { requireAuth } from '~/server/utils/auth'
import { readArticles, writeArticles } from '~/server/utils/articles'

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

  const body = await readBody<{
    title?: string
    summary?: string
    content?: Record<string, any>
  }>(event)

  // 至少需要提供一个更新字段
  if (!body?.title && !body?.summary && !body?.content) {
    throw createError({
      statusCode: 400,
      statusMessage: '至少需要提供 title、summary 或 content 中的一个字段',
    })
  }

  const articles = await readArticles()
  const index = articles.findIndex((a) => a.id === id)

  if (index === -1) {
    throw createError({
      statusCode: 404,
      statusMessage: '文章不存在',
    })
  }

  const article = articles[index]

  // 更新字段（仅更新传入的非空字段）
  if (body.title !== undefined && typeof body.title === 'string' && body.title.trim() !== '') {
    article.title = body.title.trim()
  }
  if (body.summary !== undefined) {
    article.summary = typeof body.summary === 'string' ? body.summary.trim() : article.summary
  }
  if (body.content !== undefined && typeof body.content === 'object') {
    article.content = body.content
  }

  // 自动刷新 updatedAt
  article.updatedAt = new Date().toISOString()

  articles[index] = article
  await writeArticles(articles)

  return {
    id: article.id,
    title: article.title,
    summary: article.summary,
    content: article.content,
    createdAt: article.createdAt,
    updatedAt: article.updatedAt,
  }
})
