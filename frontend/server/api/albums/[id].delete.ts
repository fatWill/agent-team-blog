import { requireAuth } from '~/server/utils/auth'
import { deleteAlbum } from '~/server/utils/albums'

/**
 * DELETE /api/albums/:id
 * 删除相册集及其所有照片（需鉴权）
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

  const deleted = deleteAlbum(id)

  if (!deleted) {
    throw createError({
      statusCode: 404,
      statusMessage: '相册不存在',
    })
  }

  return { ok: true }
})
