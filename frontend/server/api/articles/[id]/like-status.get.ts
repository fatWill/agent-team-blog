import { getPool } from '~/server/utils/db'
import type { RowDataPacket } from 'mysql2/promise'

/**
 * GET /api/articles/:id/like-status
 * 查询当前设备对文章的点赞状态
 * Query: ?deviceId=xxx
 */
export default defineEventHandler(async (event) => {
  const articleId = event.context.params?.id
  if (!articleId) {
    throw createError({ statusCode: 400, statusMessage: '无效的文章 ID' })
  }

  const query = getQuery(event)
  const deviceId = typeof query.deviceId === 'string' ? query.deviceId.trim() : ''

  const pool = getPool()

  // 查询文章点赞数
  const [articleRows] = await pool.execute(
    'SELECT like_count FROM articles WHERE id = ? LIMIT 1',
    [articleId],
  ) as [RowDataPacket[], any]

  if (articleRows.length === 0) {
    throw createError({ statusCode: 404, statusMessage: '文章不存在' })
  }

  const likeCount = Number(articleRows[0].like_count ?? 0)

  // 查询当前设备是否已点赞
  let liked = false
  if (deviceId) {
    const [likedRows] = await pool.execute(
      'SELECT 1 FROM article_likes WHERE article_id = ? AND device_id = ? LIMIT 1',
      [articleId, deviceId],
    ) as [RowDataPacket[], any]
    liked = likedRows.length > 0
  }

  return { liked, likeCount }
})
