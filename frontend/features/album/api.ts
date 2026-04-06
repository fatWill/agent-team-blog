import type {
  AlbumItem,
  AlbumListResponse,
  PhotoItem,
  PhotoListResponse,
} from './types'

/** 获取相册集列表 */
export async function apiGetAlbums(): Promise<AlbumListResponse> {
  const res = await $fetch<AlbumListResponse>('/api/albums')
  return res
}

/** 创建相册集 */
export async function apiCreateAlbum(data: { name: string; description?: string; password?: string }): Promise<AlbumItem> {
  const res = await $fetch<AlbumItem>('/api/albums', {
    method: 'POST',
    body: data,
  })
  return res
}

/** 更新相册集 */
export async function apiUpdateAlbum(id: number, data: { name?: string; description?: string; password?: string }): Promise<AlbumItem> {
  const res = await $fetch<AlbumItem>(`/api/albums/${id}`, {
    method: 'PUT',
    body: data,
  })
  return res
}

/** 删除相册集 */
export async function apiDeleteAlbum(id: number): Promise<{ ok: boolean }> {
  const res = await $fetch<{ ok: boolean }>(`/api/albums/${id}`, {
    method: 'DELETE',
  })
  return res
}

/** 获取相册下所有照片 */
export async function apiGetPhotos(albumId: number): Promise<PhotoListResponse> {
  const res = await $fetch<PhotoListResponse>(`/api/albums/${albumId}/photos`)
  return res
}

/** 向相册添加照片/视频 */
export async function apiAddPhoto(albumId: number, data: {
  url: string
  caption?: string
  mediaType?: 'image' | 'video'
  thumbnailUrl?: string
  duration?: number
}): Promise<PhotoItem> {
  const res = await $fetch<PhotoItem>(`/api/albums/${albumId}/photos`, {
    method: 'POST',
    body: data,
  })
  return res
}

/** 删除单张照片 */
export async function apiDeletePhoto(id: number): Promise<{ ok: boolean }> {
  const res = await $fetch<{ ok: boolean }>(`/api/photos/${id}`, {
    method: 'DELETE',
  })
  return res
}

/** 验证相册集密码 */
export async function apiVerifyAlbumPassword(albumId: number, password: string): Promise<{ success: boolean }> {
  const res = await $fetch<{ success: boolean }>(`/api/albums/${albumId}/verify-password`, {
    method: 'POST',
    body: { password },
  })
  return res
}

/** 验证照片密码 */
export async function apiVerifyPhotoPassword(photoId: number, password: string): Promise<{ success: boolean }> {
  const res = await $fetch<{ success: boolean }>(`/api/photos/${photoId}/verify-password`, {
    method: 'POST',
    body: { password },
  })
  return res
}

/** 更新照片信息（caption、password） */
export async function apiUpdatePhoto(id: number, data: { caption?: string; password?: string }): Promise<PhotoItem> {
  const res = await $fetch<PhotoItem>(`/api/photos/${id}`, {
    method: 'PUT',
    body: data,
  })
  return res
}
