import { getPool } from '~/server/utils/db'
import type { RowDataPacket, ResultSetHeader } from 'mysql2/promise'

/**
 * POST /api/articles/:id/like
 * 文章点赞/取消点赞（切换）
 * Body: { deviceId: string }
 */
export default defineEventHandler(async (event) => {
  const articleId = event.context.params?.id
  if (!articleId) {
    throw createError({ statusCode: 400, statusMessage: '无效的文章 ID' })
  }

  const body = await readBody<{ deviceId: string }>(event)
  const deviceId = body?.deviceId?.trim()

  if (!deviceId || deviceId.length < 8 || deviceId.length > 64) {
    throw createError({ statusCode: 400, statusMessage: '无效的设备 ID' })
  }

  const pool = getPool()

  // 检查文章是否存在
  const [articleRows] = await pool.execute(
    'SELECT id FROM articles WHERE id = ? LIMIT 1',
    [articleId],
  ) as [RowDataPacket[], any]

  if (articleRows.length === 0) {
    throw createError({ statusCode: 404, statusMessage: '文章不存在' })
  }

  // 检查是否已点赞
  const [existRows] = await pool.execute(
    'SELECT id FROM article_likes WHERE article_id = ? AND device_id = ? LIMIT 1',
    [articleId, deviceId],
  ) as [RowDataPacket[], any]

  let liked: boolean

  if (existRows.length > 0) {
    // 已点赞 → 取消点赞
    await pool.execute(
      'DELETE FROM article_likes WHERE article_id = ? AND device_id = ?',
      [articleId, deviceId],
    )
    await pool.execute(
      'UPDATE articles SET like_count = GREATEST(like_count - 1, 0) WHERE id = ?',
      [articleId],
    )
    liked = false
  }
  else {
    // 未点赞 → 点赞
    await pool.execute(
      'INSERT INTO article_likes (article_id, device_id) VALUES (?, ?)',
      [articleId, deviceId],
    )
    await pool.execute(
      'UPDATE articles SET like_count = like_count + 1 WHERE id = ?',
      [articleId],
    )
    liked = true
  }

  // 查询最新计数
  const [countRows] = await pool.execute(
    'SELECT like_count FROM articles WHERE id = ?',
    [articleId],
  ) as [RowDataPacket[], any]

  const likeCount = Number(countRows[0]?.like_count ?? 0)

  return { liked, likeCount }
})
