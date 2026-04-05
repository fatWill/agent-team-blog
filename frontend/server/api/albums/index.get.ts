import { getAlbumList } from '~/server/utils/albums'

/**
 * GET /api/albums
 * 获取所有相册集列表（公开接口，无需鉴权）
 */
export default defineEventHandler(() => {
  const list = getAlbumList()
  return { list }
})
