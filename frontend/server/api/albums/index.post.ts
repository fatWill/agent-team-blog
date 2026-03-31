import { readBody } from 'h3'
import { requireAuth } from '~/server/utils/auth'
import { createAlbum } from '~/server/utils/albums'

/**
 * POST /api/albums
 * 创建相册集（需鉴权）
 */
export default defineEventHandler(async (event) => {
  await requireAuth(event)

  const body = await readBody<{ name?: string; description?: string }>(event)

  if (!body?.name || body.name.trim() === '') {
    throw createError({
      statusCode: 400,
      statusMessage: '相册名称不能为空',
    })
  }

  const album = await createAlbum(body.name.trim(), body.description?.trim())
  return album
})
