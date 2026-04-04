// 相册领域模块 - 公共导出
export type {
  AlbumItem,
  AlbumListResponse,
  PhotoItem,
  PhotoListResponse,
} from './types'

export {
  apiGetAlbums,
  apiCreateAlbum,
  apiUpdateAlbum,
  apiDeleteAlbum,
  apiGetPhotos,
  apiAddPhoto,
  apiDeletePhoto,
  apiVerifyAlbumPassword,
  apiVerifyPhotoPassword,
  apiUpdatePhoto,
} from './api'
