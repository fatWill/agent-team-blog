import type {
  ArticleListResponse,
  ArticleDetail,
  CreateArticleRequest,
  UpdateArticleRequest,
  LoginRequest,
  LoginResponse,
  Profile,
  AlbumItem,
  AlbumListResponse,
  PhotoItem,
  PhotoListResponse,
  MessageItem,
  MessageListResponse,
} from '~/types'

/** 登录 */
export async function apiLogin(data: LoginRequest): Promise<LoginResponse> {
  const res = await $fetch<LoginResponse>('/api/auth/login', {
    method: 'POST',
    body: data,
  })
  return res
}

/** 获取文章列表 */
export async function apiFetchArticles(): Promise<ArticleListResponse> {
  const res = await $fetch<ArticleListResponse>('/api/articles')
  return res
}

/** 获取文章详情 */
export async function apiFetchArticle(id: string): Promise<ArticleDetail> {
  const res = await $fetch<ArticleDetail>(`/api/articles/${id}`)
  return res
}

/** 创建文章 */
export async function apiCreateArticle(data: CreateArticleRequest): Promise<ArticleDetail> {
  const res = await $fetch<ArticleDetail>('/api/articles', {
    method: 'POST',
    body: data,
  })
  return res
}

/** 更新文章 */
export async function apiUpdateArticle(id: string, data: UpdateArticleRequest): Promise<ArticleDetail> {
  const res = await $fetch<ArticleDetail>(`/api/articles/${id}`, {
    method: 'PUT',
    body: data,
  })
  return res
}

/** 删除文章 */
export async function apiDeleteArticle(id: string): Promise<{ ok: boolean }> {
  const res = await $fetch<{ ok: boolean }>(`/api/articles/${id}`, {
    method: 'DELETE',
  })
  return res
}

/** 获取文章列表（支持 title 关键词筛选） */
export async function apiGetArticles(params?: { title?: string }): Promise<ArticleListResponse> {
  const query: Record<string, string> = {}
  if (params?.title) {
    query.title = params.title
  }
  const res = await $fetch<ArticleListResponse>('/api/articles', {
    params: query,
  })
  return res
}

/** 上传图片 */
export async function apiUploadImage(file: File): Promise<{ url: string }> {
  const formData = new FormData()
  formData.append('file', file)
  const res = await $fetch<{ url: string }>('/api/upload', {
    method: 'POST',
    body: formData,
  })
  return res
}

/** 获取博主个人资料 */
export async function apiGetProfile(): Promise<Profile> {
  const res = await $fetch<Profile>('/api/profile')
  return res
}

/** 更新博主个人资料 */
export async function apiUpdateProfile(data: Partial<Profile>): Promise<Profile> {
  const res = await $fetch<Profile>('/api/profile', {
    method: 'PUT',
    body: data,
  })
  return res
}

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

/** 向相册添加照片 */
export async function apiAddPhoto(albumId: number, data: { url: string; caption?: string }): Promise<PhotoItem> {
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

// ====== 文章点赞 ======

/** 文章点赞/取消点赞（切换） */
export async function apiToggleArticleLike(articleId: string, deviceId: string, signal?: AbortSignal): Promise<{ liked: boolean; likeCount: number }> {
  const res = await $fetch<{ liked: boolean; likeCount: number }>(`/api/articles/${articleId}/like`, {
    method: 'POST',
    body: { deviceId },
    signal,
  })
  return res
}

/** 查询文章点赞状态 */
export async function apiGetArticleLikeStatus(articleId: string, deviceId: string): Promise<{ liked: boolean; likeCount: number }> {
  const res = await $fetch<{ liked: boolean; likeCount: number }>(`/api/articles/${articleId}/like-status`, {
    params: { deviceId },
  })
  return res
}

/** 批量查询文章点赞状态 */
export async function apiGetArticleLikeStatusBatch(ids: string[], deviceId: string): Promise<{ likedIds: string[] }> {
  const res = await $fetch<{ likedIds: string[] }>('/api/articles/like-status-batch', {
    params: { deviceId, ids: ids.join(',') },
  })
  return res
}

// ====== 留言板 ======

/** 获取留言列表 */
export async function apiGetMessages(deviceId?: string): Promise<MessageListResponse> {
  const params: Record<string, string> = {}
  if (deviceId) params.deviceId = deviceId
  const res = await $fetch<MessageListResponse>('/api/messages', { params })
  return res
}

/** 新增留言 */
export async function apiCreateMessage(data: { deviceId: string; nickname?: string; content: string }): Promise<MessageItem> {
  const res = await $fetch<MessageItem>('/api/messages', {
    method: 'POST',
    body: data,
  })
  return res
}

/** 修改留言 */
export async function apiUpdateMessage(id: number, data: { deviceId: string; nickname?: string; content: string }): Promise<MessageItem> {
  const res = await $fetch<MessageItem>(`/api/messages/${id}`, {
    method: 'PUT',
    body: data,
  })
  return res
}