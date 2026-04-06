/** 相册集列表项 */
export interface AlbumItem {
  id: number
  name: string
  description: string | null
  coverUrl: string | null
  photoCount: number
  hasPassword: boolean
  createdAt: string
  updatedAt: string
}

/** 相册集列表响应 */
export interface AlbumListResponse {
  list: AlbumItem[]
}

/** 照片列表项 */
export interface PhotoItem {
  id: number
  albumId: number
  url: string
  caption: string | null
  hasPassword: boolean
  /** 媒体类型：'image' 或 'video'，历史数据可能为 null（视为 'image'） */
  mediaType?: 'image' | 'video' | null
  /** 视频封面图 URL */
  thumbnailUrl?: string | null
  /** 视频时长（秒） */
  duration?: number | null
  createdAt: string
  updatedAt: string
  /** 点赞数（前端运行时附加，非数据库字段） */
  likes?: number
  /** 当前设备是否已点赞（前端运行时附加） */
  liked?: boolean
  /** 踩数（前端运行时附加，非数据库字段） */
  dislikes?: number
  /** 当前设备是否已踩（前端运行时附加） */
  disliked?: boolean
}

/** 照片列表响应 */
export interface PhotoListResponse {
  list: PhotoItem[]
}
