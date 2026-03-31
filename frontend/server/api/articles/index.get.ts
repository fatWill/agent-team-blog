import { getArticleList } from '~/server/utils/articles'

export default defineEventHandler(async () => {
  const list = await getArticleList()
  return { list }
})
