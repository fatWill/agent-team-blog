import { getQuery } from 'h3'
import { getArticleList } from '~/server/utils/articles'

export default defineEventHandler(async (event) => {
  const query = getQuery(event)
  const title = typeof query.title === 'string' ? query.title : undefined
  const list = await getArticleList(title)
  return { list }
})
