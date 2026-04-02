import { readBody } from 'h3'
import bcrypt from 'bcryptjs'
import { getPhotoPasswordHash } from '~/server/utils/albums'

/**
 * POST /api/photos/:id/verify-password
 * 验证照片密码（公开接口，无需鉴权）
 */
export default defineEventHandler(async (event) => {
  const photoId = Number(event.context.params?.id)
  if (!photoId || isNaN(photoId)) {
    throw createError({
      statusCode: 400,
      statusMessage: '无效的照片 ID',
    })
  }

  const body = await readBody<{ password?: string }>(event)

  if (!body?.password) {
    throw createError({
      statusCode: 400,
      statusMessage: '请输入密码',
    })
  }

  const hash = await getPhotoPasswordHash(photoId)

  // 照片不存在或未设置密码
  if (hash === null) {
    throw createError({
      statusCode: 404,
      statusMessage: '照片不存在或未设置密码',
    })
  }

  const success = await bcrypt.compare(body.password, hash)
  return { success }
})
