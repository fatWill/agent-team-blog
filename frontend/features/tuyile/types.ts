// 玩具「图一乐」领域类型定义（与 backend/docs/api/tuyile.md 对齐）

/** 单个时辰吉凶 */
export interface AlmanacHour {
  zhi: string
  ganZhi: string
  range: string
  tianShen: string
  luck: string
  chong: string
  sha: string
  yi: string[]
  ji: string[]
  jiShen: string[]
  xiongSha: string[]
}

/** 单日黄历（完整字段） */
export interface AlmanacDay {
  date: string
  week: string
  lunarYear: number
  lunarMonth: number
  lunarDay: number
  isLeapMonth: boolean
  lunarMonthCn: string
  lunarDayCn: string
  lunarDate: string
  lunarYearCn: string
  zodiacYear: string
  ganZhiYear: string
  ganZhiMonth: string
  ganZhiDay: string
  jieQi: string
  currentJieQi: string
  festivals: string[]
  yi: string[]
  ji: string[]
  chongZodiac: string
  chongGanZhi: string
  chongAge: number
  chongDesc: string
  sha: string
  naYinYear: string
  naYinMonth: string
  naYinDay: string
  zhiShen: string
  zhiShenType: string
  jianChu: string
  xiu: string
  xiuLuck: string
  xiuGong: string
  pengZu: string
  jiShen: string[]
  xiongSha: string[]
  taiShen: string
  xiShen: string
  caiShen: string
  fuShen: string
  yangGui: string
  yinGui: string
  liuYao: string
  wuHou: string
  yueXiang: string
  hours: AlmanacHour[]
  source: string
  algoFields: string[]
  mismatches: string[]
  updatedAt: string
}

/** 月度黄历单日（日历视图精简字段） */
export interface AlmanacMonthItem {
  date: string
  week: string
  lunarDate: string
  cellText: string
  ganZhiDay: string
  jieQi: string
  festivals: string[]
  zhiShenType: string
  yiTop: string[]
  jiTop: string[]
  chongZodiac: string
  source: string
}

/** 月度黄历 */
export interface AlmanacMonth {
  month: string
  year: number
  days: number
  list: AlmanacMonthItem[]
}

/** 后端统一错误体 */
export interface TuyileErrorBody {
  error: boolean
  statusCode: number
  statusMessage: string
}

// ====== AG-UI ======

export type ChatRole = 'user' | 'assistant'

/** AG-UI Message（本玩具只用纯文本） */
export interface AgUiMessage {
  id: string
  role: ChatRole
  content: string
}

/** AG-UI RunAgentInput（服务端忽略 tools/state/context） */
export interface RunAgentInput {
  threadId: string
  runId: string
  messages: AgUiMessage[]
}

export interface RunStartedEvent {
  type: 'RUN_STARTED'
  threadId?: string
  runId?: string
}

export interface TextMessageStartEvent {
  type: 'TEXT_MESSAGE_START'
  messageId: string
  role?: string
}

export interface TextMessageContentEvent {
  type: 'TEXT_MESSAGE_CONTENT'
  messageId: string
  delta: string
}

export interface TextMessageEndEvent {
  type: 'TEXT_MESSAGE_END'
  messageId: string
}

export interface RunFinishedEvent {
  type: 'RUN_FINISHED'
  threadId?: string
  runId?: string
}

export interface RunErrorEvent {
  type: 'RUN_ERROR'
  message: string
  code?: string
}

/** REASONING_* 事件：只携带状态摘要，前端不渲染成气泡、不回传 */
export interface ReasoningEvent {
  type:
    | 'REASONING_START'
    | 'REASONING_END'
    | 'REASONING_MESSAGE_START'
    | 'REASONING_MESSAGE_CONTENT'
    | 'REASONING_MESSAGE_END'
  messageId?: string
  delta?: string
}

/** CUSTOM tuyile.status 的 value */
export interface TuyileStatusValue {
  stage: 'thinking' | 'querying' | string
  text: string
  tool?: string
  range?: string
}

export interface CustomEvent {
  type: 'CUSTOM'
  name: string
  value?: unknown
}

export type AgUiEvent =
  | RunStartedEvent
  | TextMessageStartEvent
  | TextMessageContentEvent
  | TextMessageEndEvent
  | RunFinishedEvent
  | RunErrorEvent
  | ReasoningEvent
  | CustomEvent

/** 聊天气泡（前端内存态） */
export interface ChatBubble {
  id: string
  role: ChatRole
  content: string
  streaming: boolean
}

/** 思考状态条 */
export interface ChatStatus {
  stage: string
  text: string
  range?: string
}

// ====== 分享 ======

export interface ShareMessage {
  role: ChatRole
  content: string
}

export interface CreateShareResponse {
  hash: string
  url: string
}

export interface ShareDetail {
  hash: string
  messages: ShareMessage[]
  createdAt: string
}
