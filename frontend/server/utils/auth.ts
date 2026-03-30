import type { H3Event } from 'h3'
import { getCookie } from 'h3'

/** Token 存储条目 */
interface TokenEntry {
  username: string
  createdAt: number
}

/** 内存 Token 存储 Map */
const tokenStore = new Map<string, TokenEntry>()

/** Token 有效期：24 小时（毫秒） */
const TOKEN_TTL = 24 * 60 * 60 * 1000

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
 * 保存 Token 到内存 Map
 */
export function saveToken(token: string, username: string): void {
  // 清理该用户的旧 Token，避免内存泄漏
  for (const [key, entry] of tokenStore.entries()) {
    if (entry.username === username) {
      tokenStore.delete(key)
    }
  }
  tokenStore.set(token, { username, createdAt: Date.now() })
}

/**
 * 验证 Token 是否有效
 * 返回用户名或 null
 */
export function verifyToken(token: string): string | null {
  const entry = tokenStore.get(token)
  if (!entry) return null

  // 检查是否过期
  if (Date.now() - entry.createdAt > TOKEN_TTL) {
    tokenStore.delete(token)
    return null
  }

  return entry.username
}

/**
 * 从请求中提取并验证鉴权信息
 * 未通过鉴权时抛出 401 错误
 */
export function requireAuth(event: H3Event): string {
  const token = getCookie(event, 'auth_token')

  if (!token) {
    throw createError({
      statusCode: 401,
      statusMessage: '未登录，请先登录',
    })
  }

  const username = verifyToken(token)
  if (!username) {
    throw createError({
      statusCode: 401,
      statusMessage: 'Token 已过期或无效，请重新登录',
    })
  }

  return username
}
