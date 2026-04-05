import { getDb } from '~/server/utils/db'

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
  const db = getDb()

  if (action === 'dislike') {
    // INSERT OR IGNORE：重复踩静默忽略
    db.prepare(
      'INSERT OR IGNORE INTO photo_dislikes (photo_id, device_id) VALUES (?, ?)',
    ).run(photoId, deviceId)
  } else {
    // undislike：删除踩记录
    db.prepare(
      'DELETE FROM photo_dislikes WHERE photo_id = ? AND device_id = ?',
    ).run(photoId, deviceId)
  }

  // 查当前踩数
  const countRow = db.prepare(
    'SELECT COUNT(*) AS cnt FROM photo_dislikes WHERE photo_id = ?',
  ).get(photoId) as { cnt: number }
  const count = Number(countRow?.cnt ?? 0)

  // 查当前设备是否已踩
  const dislikedRow = db.prepare(
    'SELECT 1 FROM photo_dislikes WHERE photo_id = ? AND device_id = ? LIMIT 1',
  ).get(photoId, deviceId)
  const disliked = !!dislikedRow

  return { count, disliked }
})
