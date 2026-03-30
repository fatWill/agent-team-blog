import type {
  ArticleListResponse,
  ArticleDetail,
  CreateArticleRequest,
  LoginRequest,
  LoginResponse,
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
