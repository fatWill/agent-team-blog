import { getMessageList } from '~/server/utils/messages'

/**
 * GET /api/messages
 * 获取所有留言列表（公开接口）
 * Query: ?deviceId=xxx（用于判断 isOwn 和 canEdit）
 */
export default defineEventHandler(async (event) => {
  const query = getQuery(event)
  const deviceId = typeof query.deviceId === 'string' ? query.deviceId.trim() : undefined

  const list = await getMessageList(deviceId)
  return { list }
})
