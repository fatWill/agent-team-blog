/** 导航项 */
export interface NavItem {
  label: string
  to: string
}

/** 文章列表项（不含 content） */
export interface ArticleListItem {
  id: string
  title: string
  summary: string
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
  content: Record<string, unknown>
}

/** 登录请求体 */
export interface LoginRequest {
  username: string
  password: string
}

/** 登录响应 */
export interface LoginResponse {
  token: string
}

/** Tab 项 */
export interface TabItem {
  key: string
  label: string
}
