// 玩具「图一乐」领域模块 - 公共导出
export type {
  AlmanacDay,
  AlmanacHour,
  AlmanacMonth,
  AlmanacMonthItem,
  AgUiEvent,
  AgUiMessage,
  ChatBubble,
  ChatStatus,
  CreateShareResponse,
  RunAgentInput,
  ShareDetail,
  ShareMessage,
} from './types'

export {
  apiCreateTuyileShare,
  apiFetchAlmanacDay,
  apiFetchAlmanacMonth,
  apiFetchTuyileShare,
  apiStreamTuyileChat,
  readSseEvents,
  TuyileHttpError,
  TUYILE_MAX_MESSAGE_LENGTH,
  SHARE_MAX_MESSAGES,
  SHARE_MAX_CONTENT_LENGTH,
} from './api'

export { beijingToday, isValidDateStr, markdownSummary, planShare, isWechatUA } from './utils'
export { useTuyileChat } from './composables/useTuyileChat'
export { setTuyileWechatCard, rememberFirstEntryUrl } from './composables/useTuyileWechat'

export { default as TuyileShell } from './components/TuyileShell.vue'
export { default as TuyileHome } from './components/TuyileHome.vue'
export { default as AlmanacView } from './components/AlmanacView.vue'
export { default as FortuneChat } from './components/FortuneChat.vue'
export { default as ChatMessage } from './components/ChatMessage.vue'
