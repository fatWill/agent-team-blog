import type { H3Event } from 'h3'
import { getCookie, setCookie } from 'h3'
import { getRedis } from './redis'

/** Token 有效期：72 小时（秒） */
const TOKEN_TTL_SECONDS = 72 * 60 * 60

/** Redis Key 前缀 */
const TOKEN_KEY_PREFIX = 'auth_token:'

/**
 * 生成随机 Token 字符串
 */
export function generateToken(): string {
  const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'
  let token = ''
  for (let i = 0; i < 64; i++) {
    token += chars.charAt(Math.floor(Math.random() * chars.length))
  }
  return token
}

/**
 * 保存 Token 到 Redis
 * 使用 SET EX 命令，自动带 72 小时过期
 */
export async function saveToken(token: string, username: string): Promise<void> {
  const redis = getRedis()
  await redis.set(`${TOKEN_KEY_PREFIX}${token}`, username, 'EX', TOKEN_TTL_SECONDS)
}

/**
 * 删除 Token（退出登录时调用）
 */
export async function deleteToken(token: string): Promise<void> {
  const redis = getRedis()
  await redis.del(`${TOKEN_KEY_PREFIX}${token}`)
}

/**
 * 验证 Token 是否有效，通过则滚动续期
 * 返回用户名或 null
 */
export async function verifyToken(token: string): Promise<string | null> {
  const redis = getRedis()
  const username = await redis.get(`${TOKEN_KEY_PREFIX}${token}`)

  if (!username) return null

  // 滚动续期：每次验证通过重置 TTL
  await redis.expire(`${TOKEN_KEY_PREFIX}${token}`, TOKEN_TTL_SECONDS)

  return username
}

/**
 * 从请求中提取并验证鉴权信息，验证通过自动续期 cookie
 * 未通过鉴权时抛出 401 错误
 */
export async function requireAuth(event: H3Event): Promise<string> {
  const token = getCookie(event, 'auth_token')

  if (!token) {
    throw createError({
      statusCode: 401,
      statusMessage: '未登录，请先登录',
    })
  }

  const username = await verifyToken(token)
  if (!username) {
    throw createError({
      statusCode: 401,
      statusMessage: 'Token 已过期或无效，请重新登录',
    })
  }

  // 续期 cookie（滚动 72 小时）
  setCookie(event, 'auth_token', token, {
    httpOnly: true,
    maxAge: TOKEN_TTL_SECONDS,
    path: '/',
    sameSite: 'lax',
  })

  return username
}