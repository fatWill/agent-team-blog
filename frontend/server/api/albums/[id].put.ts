import { readBody } from 'h3'
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

  const body = await readBody<{ name?: string; description?: string }>(event)

  if (body?.name !== undefined && body.name.trim() === '') {
    throw createError({
      statusCode: 400,
      statusMessage: '相册名称不能为空',
    })
  }

  const updates: { name?: string; description?: string } = {}
  if (body?.name !== undefined) updates.name = body.name.trim()
  if (body?.description !== undefined) updates.description = body.description.trim()

  const album = await updateAlbum(id, updates)

  if (!album) {
    throw createError({
      statusCode: 404,
      statusMessage: '相册不存在',
    })
  }

  return album
})
