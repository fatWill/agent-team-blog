import { requireAuth } from '~/server/utils/auth'
import { deletePhoto, updateAlbumCover } from '~/server/utils/albums'

/**
 * DELETE /api/photos/:id
 * 删除单张照片（需鉴权）
 * 删除后自动更新所属相册封面
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

  const { deleted, albumId } = await deletePhoto(id)

  if (!deleted) {
    throw createError({
      statusCode: 404,
      statusMessage: '照片不存在',
    })
  }

  // 更新所属相册封面
  if (albumId) {
    await updateAlbumCover(albumId)
  }

  return { ok: true }
})
