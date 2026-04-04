/** 留言板留言项 */
export interface MessageItem {
  id: number
  nickname: string
  content: string
  isOwn: boolean
  canEdit: boolean
  createdAt: string
  updatedAt: string
}

/** 留言板列表响应 */
export interface MessageListResponse {
  list: MessageItem[]
}
