import type { H3Event } from 'h3'
import type { RowDataPacket, ResultSetHeader } from 'mysql2/promise'
import { getCookie, setCookie } from 'h3'
import { getPool } from './db'

/** Token 有效期：72 小时（毫秒） */
const TOKEN_TTL = 72 * 60 * 60 * 1000

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
 * 保存 Token 到 MySQL
 * 先删除该用户的旧 Token，再插入新 Token
 * 同时清理全局过期 Token，防止表无限增长
 */
export async function saveToken(token: string, username: string): Promise<void> {
  const pool = getPool()
  const now = Date.now()

  // 删除该用户的旧 Token + 清理全局过期 Token
  await pool.query<ResultSetHeader>(
    'DELETE FROM auth_tokens WHERE username = ? OR last_active_at < ?',
    [username, now - TOKEN_TTL],
  )

  // 插入新 Token
  await pool.query<ResultSetHeader>(
    'INSERT INTO auth_tokens (token, username, last_active_at) VALUES (?, ?, ?)',
    [token, username, now],
  )
}

/**
 * 验证 Token 是否有效，通过则滚动续期
 * 返回用户名或 null
 */
export async function verifyToken(token: string): Promise<string | null> {
  const pool = getPool()
  const now = Date.now()

  interface TokenRow extends RowDataPacket {
    username: string
    last_active_at: number
  }

  const [rows] = await pool.query<TokenRow[]>(
    'SELECT username, last_active_at FROM auth_tokens WHERE token = ?',
    [token],
  )

  if (rows.length === 0) return null

  const entry = rows[0]

  // 检查是否过期
  if (now - Number(entry.last_active_at) > TOKEN_TTL) {
    // 过期则删除
    await pool.query<ResultSetHeader>(
      'DELETE FROM auth_tokens WHERE token = ?',
      [token],
    )
    return null
  }

  // 滚动续期：刷新最后活跃时间
  await pool.query<ResultSetHeader>(
    'UPDATE auth_tokens SET last_active_at = ? WHERE token = ?',
    [now, token],
  )

  return entry.username
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
    maxAge: 72 * 60 * 60,
    path: '/',
    sameSite: 'lax',
  })

  return username
}
