import { getChangelogList } from '~/server/utils/changelog'

export default defineEventHandler(async () => {
  const changelog = await getChangelogList()
  return { changelog }
})
