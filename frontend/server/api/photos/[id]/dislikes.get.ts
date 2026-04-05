import { getDb } from '~/server/utils/db'

/**
 * GET /api/photos/:id/dislikes
 * 获取指定照片的踩数，以及当前设备是否已踩
 * query: deviceId
 */
export default defineEventHandler((event) => {
  const photoId = Number(event.context.params?.id)
  if (!photoId || isNaN(photoId)) {
    throw createError({ statusCode: 400, statusMessage: '无效的照片 ID' })
  }

  const deviceId = getQuery(event).deviceId as string | undefined
  const db = getDb()

  // 查询总踩数
  const countRow = db.prepare(
    'SELECT COUNT(*) AS cnt FROM photo_dislikes WHERE photo_id = ?',
  ).get(photoId) as { cnt: number }
  const count = Number(countRow?.cnt ?? 0)

  // 查询当前设备是否已踩
  let disliked = false
  if (deviceId) {
    const dislikedRow = db.prepare(
      'SELECT 1 FROM photo_dislikes WHERE photo_id = ? AND device_id = ? LIMIT 1',
    ).get(photoId, deviceId)
    disliked = !!dislikedRow
  }

  return { count, disliked }
})
