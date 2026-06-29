/** 微信同步状态 */
export type WechatSyncStatus = 'pending' | 'syncing' | 'success' | 'failed'

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
  /** 微信草稿 media_id */
  wechatDraftMediaId?: string
  /** 最后同步时间 */
  wechatSyncedAt?: string
  /** 同步状态 */
  wechatSyncStatus?: WechatSyncStatus
  /** 同步失败原因 */
  wechatSyncError?: string
  /** 是否自动同步 */
  wechatAutoSync?: boolean
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
  wechatAutoSync?: boolean
}

/** 微信同步日志条目 */
export interface WechatSyncLog {
  id: string
  articleId: string
  status: WechatSyncStatus
  mediaId?: string
  error?: string
  createdAt: string
}

/** 微信 access_token 状态 */
export interface WechatTokenStatus {
  cached: boolean
  remainingSeconds?: number
  lastRefreshedAt?: string
}

/** 微信服务器 IP 响应 */
export interface WechatServerIp {
  ip: string
}
