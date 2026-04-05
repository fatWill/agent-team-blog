import { readBody } from 'h3'
import bcrypt from 'bcryptjs'
import { requireAuth } from '~/server/utils/auth'
import { updateAlbum } from '~/server/utils/albums'

/**
 * PUT /api/albums/:id
 * 更新相册集（需鉴权）
 */
export default defineEventHandler(async (event) => {
  await requireAuth(event)

  const id = Number(event.context.params?.id)
  if (!id || isNaN(id)) {
    throw createError({
      statusCode: 400,
      statusMessage: '无效的相册 ID',
    })
  }

  const body = await readBody<{ name?: string; description?: string; password?: string }>(event)

  if (body?.name !== undefined && body.name.trim() === '') {
    throw createError({
      statusCode: 400,
      statusMessage: '相册名称不能为空',
    })
  }

  const updates: { name?: string; description?: string; passwordHash?: string | null } = {}
  if (body?.name !== undefined) updates.name = body.name.trim()
  if (body?.description !== undefined) updates.description = body.description.trim()

  // password 字段处理：空字符串 = 清除密码，非空 = bcrypt 哈希
  if (body?.password !== undefined) {
    if (body.password === '') {
      updates.passwordHash = null
    }
    else {
      updates.passwordHash = await bcrypt.hash(body.password.trim(), 10)
    }
  }

  const album = updateAlbum(id, updates)

  if (!album) {
    throw createError({
      statusCode: 404,
      statusMessage: '相册不存在',
    })
  }

  return album
})
