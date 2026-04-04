import type {
  ArticleListResponse,
  ArticleDetail,
  CreateArticleRequest,
  UpdateArticleRequest,
} from './types'

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
