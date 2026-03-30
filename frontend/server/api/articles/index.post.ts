import { v4 as uuidv4 } from 'uuid'
import { readBody } from 'h3'
import { requireAuth } from '~/server/utils/auth'
import { readArticles, writeArticles } from '~/server/utils/articles'
import type { Article } from '~/server/utils/articles'

export default defineEventHandler(async (event) => {
  // 鉴权校验
  requireAuth(event)

  const body = await readBody<{
    title?: string
    summary?: string
    content?: Record<string, any>
  }>(event)

  // 参数校验
  if (!body?.title || typeof body.title !== 'string' || body.title.trim() === '') {
    throw createError({
      statusCode: 400,
      statusMessage: '文章标题不能为空',
    })
  }

  if (!body?.content || typeof body.content !== 'object') {
    throw createError({
      statusCode: 400,
      statusMessage: '文章内容不能为空，且必须为 Tiptap JSON 对象',
    })
  }

  const now = new Date().toISOString()

  const newArticle: Article = {
    id: uuidv4(),
    title: body.title.trim(),
    summary: body.summary?.trim() || '',
    content: body.content,
    createdAt: now,
    updatedAt: now,
  }

  // 读取现有文章，追加新文章，写回文件
  const articles = await readArticles()
  articles.push(newArticle)
  await writeArticles(articles)

  return newArticle
})
