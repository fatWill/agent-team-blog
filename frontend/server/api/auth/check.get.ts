import { requireAuth } from '~/server/utils/auth'

/**
 * GET /api/auth/check
 * 检查当前 cookie token 是否有效，同时触发续期
 */
export default defineEventHandler((event) => {
  const username = requireAuth(event)
  return { ok: true, username }
})
