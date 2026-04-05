import { readBody } from 'h3'
import bcrypt from 'bcryptjs'
import { requireAuth } from '~/server/utils/auth'
import { updatePhoto } from '~/server/utils/albums'

/**
 * PUT /api/photos/:id
 * 更新照片信息（需鉴权）
 * 支持修改 caption 和 password
 */
export default defineEventHandler(async (event) => {
  await requireAuth(event)

  const id = Number(event.context.params?.id)
  if (!id || isNaN(id)) {
    throw createError({
      statusCode: 400,
      statusMessage: '无效的照片 ID',
    })
  }

  const body = await readBody<{ caption?: string; password?: string }>(event)

  const updates: { caption?: string; passwordHash?: string | null } = {}

  if (body?.caption !== undefined) {
    updates.caption = body.caption.trim()
  }

  // password 字段处理：空字符串 = 清除密码，非空 = bcrypt 哈希
  if (body?.password !== undefined) {
    if (body.password === '') {
      updates.passwordHash = null
    }
    else {
      updates.passwordHash = await bcrypt.hash(body.password.trim(), 10)
    }
  }

  if (Object.keys(updates).length === 0) {
    throw createError({
      statusCode: 400,
      statusMessage: '没有需要更新的字段',
    })
  }

  const photo = updatePhoto(id, updates)

  if (!photo) {
    throw createError({
      statusCode: 404,
      statusMessage: '照片不存在',
    })
  }

  return photo
})
