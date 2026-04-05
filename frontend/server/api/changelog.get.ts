import { getChangelogList } from '~/server/utils/changelog'

export default defineEventHandler(() => {
  const changelog = getChangelogList()
  return { changelog }
})
