/**
 * GET /api/theme
 * 从 Redis 读取当前用户的主题偏好
 * 通过 cookie 中的 userId 标识用户（匿名用户用 IP 作 key）
 */
export default defineEventHandler(async (event) => {
  const redis = getRedis()

  // 用 cookie 里的 uid，没有就用 IP
  const uid = getCookie(event, 'blog_uid') || getRequestIP(event) || 'anonymous'
  const key = `theme:${uid}`

  try {
    const theme = await redis.get(key)
    return { theme: theme === 'dark' ? 'dark' : 'light' }
  } catch {
    return { theme: 'light' }
  }
})
