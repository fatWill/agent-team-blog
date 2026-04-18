/** 装修文章详情 */
export interface RenovationArticle {
  id: number
  title: string
  content: string // Markdown 格式内容
  createdAt: string
  updatedAt: string
}

/** 装修文章列表项（不含 content） */
export interface RenovationArticleListItem {
  id: number
  title: string
  createdAt: string
  updatedAt: string
}
