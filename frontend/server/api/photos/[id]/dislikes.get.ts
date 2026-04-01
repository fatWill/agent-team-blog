import { getPool } from '~/server/utils/db'

/**
 * GET /api/photos/:id/dislikes
 * 获取指定照片的踩数，以及当前设备是否已踩
 * query: deviceId
 */
export default defineEventHandler(async (event) => {
  const photoId = Number(event.context.params?.id)
  if (!photoId || isNaN(photoId)) {
    throw createError({ statusCode: 400, statusMessage: '无效的照片 ID' })
  }

  const deviceId = getQuery(event).deviceId as string | undefined
  const pool = getPool()

  // 查询总踩数
  const [countRows] = await pool.execute(
    'SELECT COUNT(*) AS cnt FROM photo_dislikes WHERE photo_id = ?',
    [photoId]
  ) as any[]
  const count = Number(countRows[0]?.cnt ?? 0)

  // 查询当前设备是否已踩
  let disliked = false
  if (deviceId) {
    const [dislikedRows] = await pool.execute(
      'SELECT 1 FROM photo_dislikes WHERE photo_id = ? AND device_id = ? LIMIT 1',
      [photoId, deviceId]
    ) as any[]
    disliked = (dislikedRows as any[]).length > 0
  }

  return { count, disliked }
})
