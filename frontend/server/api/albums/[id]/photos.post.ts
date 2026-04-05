import { readBody } from 'h3'
import bcrypt from 'bcryptjs'
import { requireAuth } from '~/server/utils/auth'
import { getAlbumById, addPhoto, updateAlbumCover } from '~/server/utils/albums'

/**
 * POST /api/albums/:id/photos
 * 向相册添加照片（需鉴权）
 * 添加后自动更新相册封面
 */
export default defineEventHandler(async (event) => {
  await requireAuth(event)

  const albumId = Number(event.context.params?.id)
  if (!albumId || isNaN(albumId)) {
    throw createError({
      statusCode: 400,
      statusMessage: '无效的相册 ID',
    })
  }

  // 检查相册是否存在
  const album = getAlbumById(albumId)
  if (!album) {
    throw createError({
      statusCode: 404,
      statusMessage: '相册不存在',
    })
  }

  const body = await readBody<{ url?: string; caption?: string; password?: string }>(event)

  if (!body?.url || body.url.trim() === '') {
    throw createError({
      statusCode: 400,
      statusMessage: '图片 URL 不能为空',
    })
  }

  // 如果设置了密码，bcrypt 哈希后存储
  let passwordHash: string | null = null
  if (body.password && body.password.trim() !== '') {
    passwordHash = await bcrypt.hash(body.password.trim(), 10)
  }

  const photo = addPhoto(albumId, body.url.trim(), body.caption?.trim(), passwordHash)

  // 更新相册封面
  updateAlbumCover(albumId)

  return photo
})
