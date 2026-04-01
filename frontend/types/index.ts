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
  coverImage: string
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

/** 单个版本更新日志 */
export interface ChangelogItem {
  version: string
  date: string
  logs: string[]
  createdAt: string
  updatedAt: string
}

/** 更新日志接口响应 */
export interface ChangelogResponse {
  changelog: ChangelogItem[]
}

/** 博主个人资料 */
export interface Profile {
  avatar: string
  bio: string
}

/** 相册集列表项 */
export interface AlbumItem {
  id: number
  name: string
  description: string | null
  coverUrl: string | null
  photoCount: number
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