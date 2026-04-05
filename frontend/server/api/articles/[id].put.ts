import { readBody } from 'h3'
import { requireAuth } from '~/server/utils/auth'
import { updateArticle } from '~/server/utils/articles'

export default defineEventHandler(async (event) => {
  // 鉴权校验
  await requireAuth(event)

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
    coverImage?: string
    content?: Record<string, any>
  }>(event)

  // 至少需要提供一个更新字段
  if (!body?.title && body?.summary === undefined && body?.coverImage === undefined && !body?.content) {
    throw createError({
      statusCode: 400,
      statusMessage: '至少需要提供 title、summary、coverImage 或 content 中的一个字段',
    })
  }

  // 构建更新对象
  const updates: { title?: string; summary?: string; coverImage?: string; content?: Record<string, any> } = {}

  if (body.title !== undefined && typeof body.title === 'string' && body.title.trim() !== '') {
    updates.title = body.title.trim()
  }
  if (body.summary !== undefined) {
    updates.summary = typeof body.summary === 'string' ? body.summary.trim() : undefined
  }
  if (body.coverImage !== undefined) {
    updates.coverImage = typeof body.coverImage === 'string' ? body.coverImage.trim() : ''
  }
  if (body.content !== undefined && typeof body.content === 'object') {
    updates.content = body.content
  }

  const article = updateArticle(id, updates)

  if (!article) {
    throw createError({
      statusCode: 404,
      statusMessage: '文章不存在',
    })
  }

  return article
})
