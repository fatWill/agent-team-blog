import { getCookie, deleteCookie } from 'h3'
import { deleteToken } from '~/server/utils/auth'

/**
 * POST /api/auth/logout
 * 退出登录：删除 Redis token + 清除 cookie
 */
export default defineEventHandler(async (event) => {
  const token = getCookie(event, 'auth_token')

  if (token) {
    // 从 Redis 删除 token
    await deleteToken(token)
  }

  // 清除客户端 cookie
  deleteCookie(event, 'auth_token', {
    httpOnly: true,
    path: '/',
    sameSite: 'lax',
  })

  return { ok: true }
})
