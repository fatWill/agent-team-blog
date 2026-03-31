import { readBody, setCookie } from 'h3'
import { generateToken, saveToken } from '~/server/utils/auth'

/** 硬编码的账号信息 */
const VALID_USERNAME = 'fafa'
const VALID_PASSWORD = '***REDACTED_ADMIN_PWD***'

export default defineEventHandler(async (event) => {
  const body = await readBody<{ username?: string; password?: string }>(event)

  // 参数校验
  if (!body?.username || !body?.password) {
    throw createError({
      statusCode: 400,
      statusMessage: '请提供用户名和密码',
    })
  }

  // 验证账号密码
  if (body.username !== VALID_USERNAME || body.password !== VALID_PASSWORD) {
    throw createError({
      statusCode: 401,
      statusMessage: '用户名或密码错误',
    })
  }

  // 生成 Token 并持久化到 MySQL
  const token = generateToken()
  await saveToken(token, body.username)

  // 将 Token 写入 httpOnly cookie
  setCookie(event, 'auth_token', token, {
    httpOnly: true,
    maxAge: 60 * 60 * 72, // 72 小时
    path: '/',
    sameSite: 'lax',
  })

  return {
    token,
  }
})
