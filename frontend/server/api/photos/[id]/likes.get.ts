import { getDb } from '~/server/utils/db'

/**
 * GET /api/photos/:id/likes
 * 获取指定照片的点赞数，以及当前设备是否已点赞
 */
export default defineEventHandler((event) => {
  const photoId = Number(event.context.params?.id)
  if (!photoId || isNaN(photoId)) {
    throw createError({ statusCode: 400, statusMessage: '无效的照片 ID' })
  }

  const deviceId = getQuery(event).deviceId as string | undefined
  const db = getDb()

  // 查询总点赞数
  const countRow = db.prepare(
    'SELECT COUNT(*) AS cnt FROM photo_likes WHERE photo_id = ?',
  ).get(photoId) as { cnt: number }
  const count = Number(countRow?.cnt ?? 0)

  // 查询当前设备是否已点赞
  let liked = false
  if (deviceId) {
    const likedRow = db.prepare(
      'SELECT 1 FROM photo_likes WHERE photo_id = ? AND device_id = ? LIMIT 1',
    ).get(photoId, deviceId)
    liked = !!likedRow
  }

  return { count, liked }
})
