/** 文章列表项（不含 content） */
export interface ArticleListItem {
  id: string
  title: string
  summary: string
  coverImage: string
  likeCount: number
  views: number
  createdAt: string
  updatedAt: string
}

/** 文章详情（含 Tiptap content） */
export interface ArticleDetail extends ArticleListItem {
  content: Record<string, unknown>
}

/** 文章列表接口响应 */
export interface ArticleListResponse {
  list: ArticleListItem[]
}

/** 创建文章请求体 */
export interface CreateArticleRequest {
  title: string
  summary?: string
  coverImage?: string
  content: Record<string, unknown>
}

/** 更新文章请求体（所有字段可选） */
export interface UpdateArticleRequest {
  title?: string
  summary?: string
  coverImage?: string
  content?: Record<string, unknown>
}
