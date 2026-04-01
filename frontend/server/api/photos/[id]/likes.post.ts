import { getPool } from '~/server/utils/db'

/**
 * POST /api/photos/:id/likes
 * 给指定照片点赞，每个设备只能点赞一次（重复点赞直接返回当前状态，不报错）
 */
export default defineEventHandler(async (event) => {
  const photoId = Number(event.context.params?.id)
  if (!photoId || isNaN(photoId)) {
    throw createError({ statusCode: 400, statusMessage: '无效的照片 ID' })
  }

  const body = await readBody<{ deviceId: string }>(event)
  const deviceId = body?.deviceId?.trim()

  if (!deviceId || deviceId.length < 8 || deviceId.length > 64) {
    throw createError({ statusCode: 400, statusMessage: '无效的设备 ID' })
  }

  const pool = getPool()

  // INSERT IGNORE：重复点赞静默忽略（不报错）
  await pool.execute(
    'INSERT IGNORE INTO photo_likes (photo_id, device_id) VALUES (?, ?)',
    [photoId, deviceId]
  )

  // 返回最新点赞数
  const [countRows] = await pool.execute(
    'SELECT COUNT(*) AS cnt FROM photo_likes WHERE photo_id = ?',
    [photoId]
  ) as any[]
  const count = Number(countRows[0]?.cnt ?? 0)

  return { count, liked: true }
})
