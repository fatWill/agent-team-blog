// 留言板领域模块 - 公共导出
export type {
  MessageItem,
  MessageListResponse,
} from './types'

export {
  apiGetMessages,
  apiCreateMessage,
  apiDeleteMessage,
} from './api'
