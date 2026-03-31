import { getAlbumById, getPhotosByAlbumId } from '~/server/utils/albums'

/**
 * GET /api/albums/:id/photos
 * 获取指定相册集的所有照片（公开接口，无需鉴权）
 */
export default defineEventHandler(async (event) => {
  const albumId = Number(event.context.params?.id)
  if (!albumId || isNaN(albumId)) {
    throw createError({
      statusCode: 400,
      statusMessage: '无效的相册 ID',
    })
  }

  // 检查相册是否存在
  const album = await getAlbumById(albumId)
  if (!album) {
    throw createError({
      statusCode: 404,
      statusMessage: '相册不存在',
    })
  }

  const list = await getPhotosByAlbumId(albumId)
  return { list }
})
