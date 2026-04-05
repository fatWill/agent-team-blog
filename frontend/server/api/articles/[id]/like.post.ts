import { getDb } from '~/server/utils/db'

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

  const db = getDb()

  // 检查文章是否存在
  const articleRow = db.prepare(
    'SELECT id FROM articles WHERE id = ? LIMIT 1',
  ).get(articleId)

  if (!articleRow) {
    throw createError({ statusCode: 404, statusMessage: '文章不存在' })
  }

  // 检查是否已点赞
  const existRow = db.prepare(
    'SELECT id FROM article_likes WHERE article_id = ? AND device_id = ? LIMIT 1',
  ).get(articleId, deviceId)

  let liked: boolean

  if (existRow) {
    // 已点赞 → 取消点赞
    db.prepare('DELETE FROM article_likes WHERE article_id = ? AND device_id = ?').run(articleId, deviceId)
    db.prepare('UPDATE articles SET like_count = MAX(like_count - 1, 0) WHERE id = ?').run(articleId)
    liked = false
  }
  else {
    // 未点赞 → 点赞
    db.prepare('INSERT INTO article_likes (article_id, device_id) VALUES (?, ?)').run(articleId, deviceId)
    db.prepare('UPDATE articles SET like_count = like_count + 1 WHERE id = ?').run(articleId)
    liked = true
  }

  // 查询最新计数
  const countRow = db.prepare(
    'SELECT like_count FROM articles WHERE id = ?',
  ).get(articleId) as { like_count: number }

  const likeCount = Number(countRow?.like_count ?? 0)

  return { liked, likeCount }
})
