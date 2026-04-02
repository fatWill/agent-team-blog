import { getPool } from '~/server/utils/db'
import type { RowDataPacket } from 'mysql2/promise'

/**
 * PUT /api/messages/:id
 * 修改留言（公开接口，通过 deviceId 验证归属）
 * Body: { deviceId: string, nickname?: string, content: string }
 */
export default defineEventHandler(async (event) => {
  const id = Number(event.context.params?.id)
  if (!id || isNaN(id)) {
    throw createError({ statusCode: 400, statusMessage: '无效的留言 ID' })
  }

  const body = await readBody<{ deviceId: string; nickname?: string; content: string }>(event)

  const deviceId = body?.deviceId?.trim()
  if (!deviceId || deviceId.length < 8 || deviceId.length > 64) {
    throw createError({ statusCode: 400, statusMessage: '无效的设备 ID' })
  }

  const content = body?.content?.trim()
  if (!content) {
    throw createError({ statusCode: 400, statusMessage: '留言内容不能为空' })
  }

  const nickname = body?.nickname?.trim() || null

  const pool = getPool()

  // 查找留言
  const [rows] = await pool.execute(
    'SELECT id, device_id, last_modified_date FROM messages WHERE id = ? LIMIT 1',
    [id],
  ) as [RowDataPacket[], any]

  if (rows.length === 0) {
    throw createError({ statusCode: 404, statusMessage: '留言不存在' })
  }

  const existing = rows[0]

  // 验证设备归属
  if (existing.device_id !== deviceId) {
    throw createError({ statusCode: 403, statusMessage: '无权修改此留言' })
  }

  // 计算今天的日期字符串（UTC+8）
  const now = new Date()
  const utc8 = new Date(now.getTime() + 8 * 60 * 60 * 1000)
  const todayStr = utc8.toISOString().slice(0, 10)

  // 检查今天是否已修改过
  const rawDate = existing.last_modified_date
  let lastModDate: string | null = null
  if (rawDate instanceof Date) {
    // MySQL DATE 类型返回的 Date 对象是本地时间，直接提取年月日（避免 toISOString 时区偏移）
    const y = rawDate.getFullYear()
    const m = String(rawDate.getMonth() + 1).padStart(2, '0')
    const d = String(rawDate.getDate()).padStart(2, '0')
    lastModDate = `${y}-${m}-${d}`
  }
  else if (rawDate) {
    lastModDate = String(rawDate).slice(0, 10)
  }
  if (lastModDate === todayStr) {
    throw createError({
      statusCode: 403,
      statusMessage: '今天已经修改过了，明天再来吧',
    })
  }

  // 更新留言并记录修改日期
  await pool.execute(
    'UPDATE messages SET nickname = ?, content = ?, last_modified_date = ? WHERE id = ?',
    [nickname, content, todayStr, id],
  )

  // 查询更新后的记录
  const [updatedRows] = await pool.execute(
    'SELECT id, nickname, content, last_modified_date, created_at, updated_at FROM messages WHERE id = ?',
    [id],
  ) as [RowDataPacket[], any]

  const row = updatedRows[0]
  return {
    id: row.id,
    nickname: row.nickname || '匿名',
    content: row.content,
    isOwn: true,
    canEdit: false, // 刚修改过，今天不能再改
    createdAt: new Date(row.created_at).toISOString(),
    updatedAt: new Date(row.updated_at).toISOString(),
  }
})
