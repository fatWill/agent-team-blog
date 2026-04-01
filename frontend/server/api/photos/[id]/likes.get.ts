import { getPool } from '~/server/utils/db'

/**
 * GET /api/photos/:id/likes
 * 获取指定照片的点赞数，以及当前设备是否已点赞
 */
export default defineEventHandler(async (event) => {
  const photoId = Number(event.context.params?.id)
  if (!photoId || isNaN(photoId)) {
    throw createError({ statusCode: 400, statusMessage: '无效的照片 ID' })
  }

  const deviceId = getQuery(event).deviceId as string | undefined
  const pool = getPool()

  // 查询总点赞数
  const [countRows] = await pool.execute(
    'SELECT COUNT(*) AS cnt FROM photo_likes WHERE photo_id = ?',
    [photoId]
  ) as any[]
  const count = Number(countRows[0]?.cnt ?? 0)

  // 查询当前设备是否已点赞
  let liked = false
  if (deviceId) {
    const [likedRows] = await pool.execute(
      'SELECT 1 FROM photo_likes WHERE photo_id = ? AND device_id = ? LIMIT 1',
      [photoId, deviceId]
    ) as any[]
    liked = (likedRows as any[]).length > 0
  }

  return { count, liked }
})
