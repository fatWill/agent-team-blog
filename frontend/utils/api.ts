import type {
  ArticleListResponse,
  ArticleDetail,
  CreateArticleRequest,
  UpdateArticleRequest,
  LoginRequest,
  LoginResponse,
  Profile,
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