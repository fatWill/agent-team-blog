import { getDb } from '~/server/utils/db'

/**
 * GET /api/messages
 * 获取所有留言列表（公开接口）
 * Query: ?deviceId=xxx（可选，用于标识 isOwn 和 canEdit）
 */
export default defineEventHandler((event) => {
  const query = getQuery(event)
  const deviceId = typeof query.deviceId === 'string' ? query.deviceId.trim() : ''

  const db = getDb()

  const rows = db.prepare(
    'SELECT id, device_id, nickname, content, last_modified_date, created_at, updated_at FROM messages ORDER BY created_at DESC',
  ).all() as {
    id: number
    device_id: string
    nickname: string | null
    content: string
    last_modified_date: string | null
    created_at: string
    updated_at: string
  }[]

  // 计算今天的日期字符串（UTC+8）
  const now = new Date()
  const utc8 = new Date(now.getTime() + 8 * 60 * 60 * 1000)
  const todayStr = utc8.toISOString().slice(0, 10) // "YYYY-MM-DD"

  const list = rows.map((row) => {
    const isOwn = deviceId !== '' && row.device_id === deviceId
    const lastModDate = row.last_modified_date ? String(row.last_modified_date).slice(0, 10) : null
    const canEdit = isOwn && lastModDate !== todayStr

    return {
      id: row.id,
      nickname: row.nickname || '匿名',
      content: row.content,
      isOwn,
      canEdit,
      createdAt: row.created_at,
      updatedAt: row.updated_at,
    }
  })

  return { list }
})
