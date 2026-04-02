import { getPool } from '~/server/utils/db'
import type { RowDataPacket } from 'mysql2/promise'

/**
 * GET /api/articles/like-status-batch
 * 批量查询当前设备对多篇文章的点赞状态
 * Query: ?deviceId=xxx&ids=1,2,3,4,5
 */
export default defineEventHandler(async (event) => {
  const query = getQuery(event)

  const deviceId = typeof query.deviceId === 'string' ? query.deviceId.trim() : ''
  if (!deviceId) {
    throw createError({ statusCode: 400, statusMessage: 'deviceId 不能为空' })
  }

  const idsParam = typeof query.ids === 'string' ? query.ids.trim() : ''
  if (!idsParam) {
    return { likedIds: [] }
  }

  // 解析并过滤有效的文章 ID，最多取前 100 个
  const ids = idsParam
    .split(',')
    .map(s => s.trim())
    .filter(s => s.length > 0)
    .slice(0, 100)

  if (ids.length === 0) {
    return { likedIds: [] }
  }

  const pool = getPool()

  // 构建 IN 查询的占位符
  const placeholders = ids.map(() => '?').join(',')
  const params = [...ids, deviceId]

  const [rows] = await pool.execute(
    `SELECT article_id FROM article_likes WHERE article_id IN (${placeholders}) AND device_id = ?`,
    params,
  ) as [RowDataPacket[], any]

  const likedIds = rows.map((row: RowDataPacket) => row.article_id)

  return { likedIds }
})
