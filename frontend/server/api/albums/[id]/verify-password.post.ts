import { readBody } from 'h3'
import bcrypt from 'bcryptjs'
import { getAlbumPasswordHash } from '~/server/utils/albums'

/**
 * POST /api/albums/:id/verify-password
 * 验证相册集密码（公开接口，无需鉴权）
 */
export default defineEventHandler(async (event) => {
  const albumId = Number(event.context.params?.id)
  if (!albumId || isNaN(albumId)) {
    throw createError({
      statusCode: 400,
      statusMessage: '无效的相册 ID',
    })
  }

  const body = await readBody<{ password?: string }>(event)

  if (!body?.password) {
    throw createError({
      statusCode: 400,
      statusMessage: '请输入密码',
    })
  }

  const hash = await getAlbumPasswordHash(albumId)

  // 相册不存在或未设置密码
  if (hash === null) {
    throw createError({
      statusCode: 404,
      statusMessage: '相册不存在或未设置密码',
    })
  }

  const success = await bcrypt.compare(body.password, hash)
  return { success }
})
