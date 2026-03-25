/** 文章信息 */
export interface Article {
  id: number
  title: string
  summary: string
  content: string
  coverImage: string
  category: Category
  tags: Tag[]
  author: Author
  createdAt: string
  updatedAt: string
  readingTime: number
  viewCount: number
  likeCount: number
}

/** 文章分类 */
export interface Category {
  id: number
  name: string
  slug: string
  articleCount?: number
}

/** 文章标签 */
export interface Tag {
  id: number
  name: string
  slug: string
}

/** 作者信息 */
export interface Author {
  id: number
  name: string
  avatar: string
  bio: string
}

/** 分页响应 */
export interface PaginatedResponse<T> {
  data: T[]
  total: number
  page: number
  pageSize: number
  totalPages: number
}

/** 导航菜单项 */
export interface NavItem {
  label: string
  to: string
  icon?: string
}
