import { readArticles } from '~/server/utils/articles'

export default defineEventHandler(async () => {
  const articles = await readArticles()

  // 返回列表，不包含 content 字段，按创建时间倒序
  const list = articles
    .map(({ id, title, summary, createdAt, updatedAt }) => ({
      id,
      title,
      summary,
      createdAt,
      updatedAt,
    }))
    .sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime())

  return { list }
})
