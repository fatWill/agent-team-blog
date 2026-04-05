import { getDb } from '~/server/utils/db'

/**
 * GET /api/articles/:id/like-status
 * 查询当前设备对文章的点赞状态
 * Query: ?deviceId=xxx
 */
export default defineEventHandler((event) => {
  const articleId = event.context.params?.id
  if (!articleId) {
    throw createError({ statusCode: 400, statusMessage: '无效的文章 ID' })
  }

  const query = getQuery(event)
  const deviceId = typeof query.deviceId === 'string' ? query.deviceId.trim() : ''

  const db = getDb()

  // 查询文章点赞数
  const articleRow = db.prepare(
    'SELECT like_count FROM articles WHERE id = ? LIMIT 1',
  ).get(articleId) as { like_count: number } | undefined

  if (!articleRow) {
    throw createError({ statusCode: 404, statusMessage: '文章不存在' })
  }

  const likeCount = Number(articleRow.like_count ?? 0)

  // 查询当前设备是否已点赞
  let liked = false
  if (deviceId) {
    const likedRow = db.prepare(
      'SELECT 1 FROM article_likes WHERE article_id = ? AND device_id = ? LIMIT 1',
    ).get(articleId, deviceId)
    liked = !!likedRow
  }

  return { liked, likeCount }
})
