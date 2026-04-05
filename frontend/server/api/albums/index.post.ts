import { readBody } from 'h3'
import bcrypt from 'bcryptjs'
import { requireAuth } from '~/server/utils/auth'
import { createAlbum } from '~/server/utils/albums'

/**
 * POST /api/albums
 * 创建相册集（需鉴权）
 */
export default defineEventHandler(async (event) => {
  await requireAuth(event)

  const body = await readBody<{ name?: string; description?: string; password?: string }>(event)

  if (!body?.name || body.name.trim() === '') {
    throw createError({
      statusCode: 400,
      statusMessage: '相册名称不能为空',
    })
  }

  // 如果设置了密码，bcrypt 哈希后存储
  let passwordHash: string | null = null
  if (body.password && body.password.trim() !== '') {
    passwordHash = await bcrypt.hash(body.password.trim(), 10)
  }

  const album = createAlbum(body.name.trim(), body.description?.trim(), passwordHash)
  return album
})
