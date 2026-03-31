/**
 * POST /api/theme
 * 保存用户主题偏好到 Redis，同时写 cookie
 * Body: { theme: 'dark' | 'light' }
 */
export default defineEventHandler(async (event) => {
  const body = await readBody(event)
  const theme: 'dark' | 'light' = body?.theme === 'dark' ? 'dark' : 'light'

  const redis = getRedis()

  // 确保 uid cookie 存在（30天有效）
  let uid = getCookie(event, 'blog_uid')
  if (!uid) {
    uid = `${Date.now()}-${Math.random().toString(36).slice(2, 8)}`
    setCookie(event, 'blog_uid', uid, {
      maxAge: 60 * 60 * 24 * 30,
      httpOnly: true,
      path: '/',
    })
  }

  const key = `theme:${uid}`

  try {
    // 主题偏好存 30 天
    await redis.set(key, theme, 'EX', 60 * 60 * 24 * 30)
  } catch (err) {
    console.error('[Redis] set theme error:', err)
  }

  // 同时写一个客户端可读的 cookie，供 SSR 直接读（不需要再查 Redis）
  setCookie(event, 'color-mode', theme, {
    maxAge: 60 * 60 * 24 * 30,
    httpOnly: false,
    path: '/',
  })

  return { ok: true, theme }
})
