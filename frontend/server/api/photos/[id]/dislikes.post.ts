import { getPool } from '~/server/utils/db'

/**
 * POST /api/photos/:id/dislikes
 * 给指定照片踩/取消踩（用 device_id 防重复，可切换）
 * body: { deviceId: string, action: 'dislike' | 'undislike' }
 */
export default defineEventHandler(async (event) => {
  const photoId = Number(event.context.params?.id)
  if (!photoId || isNaN(photoId)) {
    throw createError({ statusCode: 400, statusMessage: '无效的照片 ID' })
  }

  const body = await readBody<{ deviceId: string; action?: 'dislike' | 'undislike' }>(event)
  const deviceId = body?.deviceId?.trim()

  if (!deviceId || deviceId.length < 8 || deviceId.length > 64) {
    throw createError({ statusCode: 400, statusMessage: '无效的设备 ID' })
  }

  const action = body?.action ?? 'dislike'
  const pool = getPool()

  if (action === 'dislike') {
    // INSERT IGNORE：重复踩静默忽略
    await pool.execute(
      'INSERT IGNORE INTO photo_dislikes (photo_id, device_id) VALUES (?, ?)',
      [photoId, deviceId]
    )
  } else {
    // undislike：删除踩记录
    await pool.execute(
      'DELETE FROM photo_dislikes WHERE photo_id = ? AND device_id = ?',
      [photoId, deviceId]
    )
  }

  // 查当前踩数
  const [countRows] = await pool.execute(
    'SELECT COUNT(*) AS cnt FROM photo_dislikes WHERE photo_id = ?',
    [photoId]
  ) as any[]
  const count = Number(countRows[0]?.cnt ?? 0)

  // 查当前设备是否已踩
  const [dislikedRows] = await pool.execute(
    'SELECT 1 FROM photo_dislikes WHERE photo_id = ? AND device_id = ? LIMIT 1',
    [photoId, deviceId]
  ) as any[]
  const disliked = (dislikedRows as any[]).length > 0

  return { count, disliked }
})
