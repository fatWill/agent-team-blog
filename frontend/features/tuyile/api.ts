import type { AgUiEvent, AlmanacDay, AlmanacMonth, CreateShareResponse, RunAgentInput, ShareDetail, ShareMessage } from './types'

export const TUYILE_API_BASE = '/api/toys/tuyile'
/** 单条消息最大字数（与后端 MaxMessageRunes 一致） */
export const TUYILE_MAX_MESSAGE_LENGTH = 500
export const ALMANAC_MIN_DATE = '1900-01-01'
export const ALMANAC_MAX_DATE = '2100-12-31'

/** 单日黄历；不传 date 时后端取北京时间今天 */
export function apiFetchAlmanacDay(date?: string): Promise<AlmanacDay> {
  return $fetch<AlmanacDay>(`${TUYILE_API_BASE}/almanac`, {
    query: date ? { date } : undefined,
  })
}

/** 月度黄历（YYYY-MM） */
export function apiFetchAlmanacMonth(month: string): Promise<AlmanacMonth> {
  return $fetch<AlmanacMonth>(`${TUYILE_API_BASE}/almanac/month`, {
    query: { month },
  })
}

function isRecord(v: unknown): v is Record<string, unknown> {
  return typeof v === 'object' && v !== null
}

/** 从 $fetch 错误中提取后端 statusMessage */
export function getTuyileErrorMessage(err: unknown, fallback: string): string {
  if (isRecord(err)) {
    const data = err.data
    if (isRecord(data) && typeof data.statusMessage === 'string' && data.statusMessage) {
      return data.statusMessage
    }
  }
  return fallback
}

/** 运行前拒绝（HTTP 非 2xx） */
export class TuyileHttpError extends Error {
  status: number
  retryAfter: number | null

  constructor(status: number, message: string, retryAfter: number | null) {
    super(message)
    this.name = 'TuyileHttpError'
    this.status = status
    this.retryAfter = retryAfter
  }
}

function defaultHttpMessage(status: number): string {
  if (status === 400) return '请求不太对劲，换个说法再试试～'
  if (status === 429) return '算得太勤快啦，歇一会儿再来～'
  if (status === 503) return '图一乐小摊排队中，稍后再试～'
  return '图一乐暂时联系不上，稍后再试～'
}

const EVENT_TYPES = new Set([
  'RUN_STARTED',
  'TEXT_MESSAGE_START',
  'TEXT_MESSAGE_CONTENT',
  'TEXT_MESSAGE_END',
  'RUN_FINISHED',
  'RUN_ERROR',
  'REASONING_START',
  'REASONING_END',
  'REASONING_MESSAGE_START',
  'REASONING_MESSAGE_CONTENT',
  'REASONING_MESSAGE_END',
  'CUSTOM',
])

function isAgUiEvent(v: unknown): v is AgUiEvent {
  return isRecord(v) && typeof v.type === 'string' && EVENT_TYPES.has(v.type)
}

/**
 * 解析 text/event-stream：按空行分隔事件，拼接 data 行；
 * 以 ":" 开头的注释行（如 ": keep-alive"）直接忽略；未知事件类型忽略。
 */
export async function readSseEvents(
  body: ReadableStream<Uint8Array>,
  onEvent: (event: AgUiEvent) => void,
): Promise<void> {
  const reader = body.getReader()
  const decoder = new TextDecoder()
  let buffer = ''
  let dataLines: string[] = []

  const dispatch = () => {
    if (!dataLines.length) return
    const raw = dataLines.join('\n')
    dataLines = []
    let parsed: unknown
    try {
      parsed = JSON.parse(raw)
    } catch {
      return
    }
    if (isAgUiEvent(parsed)) onEvent(parsed)
  }

  const handleLine = (line: string) => {
    if (line === '') {
      dispatch()
      return
    }
    if (line.startsWith(':')) return
    const idx = line.indexOf(':')
    const field = idx === -1 ? line : line.slice(0, idx)
    let value = idx === -1 ? '' : line.slice(idx + 1)
    if (value.startsWith(' ')) value = value.slice(1)
    if (field === 'data') dataLines.push(value)
  }

  try {
    while (true) {
      const { done, value } = await reader.read()
      if (done) break
      buffer += decoder.decode(value, { stream: true })
      let nl = buffer.indexOf('\n')
      while (nl !== -1) {
        let line = buffer.slice(0, nl)
        buffer = buffer.slice(nl + 1)
        if (line.endsWith('\r')) line = line.slice(0, -1)
        handleLine(line)
        nl = buffer.indexOf('\n')
      }
    }
    buffer += decoder.decode()
    if (buffer) handleLine(buffer.replace(/\r$/, ''))
    dispatch()
  } finally {
    reader.releaseLock()
  }
}

/**
 * AG-UI 运势聊天：POST RunAgentInput，逐个事件回调。
 * - HTTP 非 2xx 抛 TuyileHttpError（携带 Retry-After 秒数）
 * - signal 中止时抛 AbortError
 */
export async function apiStreamTuyileChat(
  input: RunAgentInput,
  onEvent: (event: AgUiEvent) => void,
  signal: AbortSignal,
): Promise<void> {
  const res = await fetch(`${TUYILE_API_BASE}/chat`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Accept: 'text/event-stream',
    },
    body: JSON.stringify(input),
    signal,
  })

  if (!res.ok) {
    const raw = res.headers.get('Retry-After')
    const retryAfter = raw && /^\d+$/.test(raw.trim()) ? Number(raw.trim()) : null
    let message = ''
    try {
      const data: unknown = await res.json()
      if (isRecord(data) && typeof data.statusMessage === 'string') message = data.statusMessage
    } catch {
      // 非 JSON 错误体，使用默认文案
    }
    throw new TuyileHttpError(res.status, message || defaultHttpMessage(res.status), retryAfter)
  }

  if (!res.body) {
    throw new TuyileHttpError(res.status, defaultHttpMessage(500), null)
  }

  await readSseEvents(res.body, onEvent)
}

// ====== 分享 ======

/** 与后端 share 接口限制一致 */
export const SHARE_MAX_MESSAGES = 50
export const SHARE_MAX_CONTENT_LENGTH = 4000
/** 请求体上限 128KB，留余量 */
export const SHARE_MAX_BODY_BYTES = 120 * 1024

/** 创建分享：失败抛 TuyileHttpError（400/413/429 带后端文案与 Retry-After） */
export async function apiCreateTuyileShare(messages: ShareMessage[]): Promise<CreateShareResponse> {
  const res = await fetch(`${TUYILE_API_BASE}/share`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', Accept: 'application/json' },
    body: JSON.stringify({ messages }),
  })
  if (!res.ok) {
    const raw = res.headers.get('Retry-After')
    const retryAfter = raw && /^\d+$/.test(raw.trim()) ? Number(raw.trim()) : null
    let message = ''
    try {
      const data: unknown = await res.json()
      if (isRecord(data) && typeof data.statusMessage === 'string') message = data.statusMessage
    } catch {
      // 非 JSON 错误体
    }
    const fallback = res.status === 413
      ? '分享内容太长啦，删减一些再试～'
      : res.status === 429 ? '分享太频繁啦，歇一会儿再来～' : '分享失败，稍后再试～'
    throw new TuyileHttpError(res.status, message || fallback, retryAfter)
  }
  const data: unknown = await res.json()
  if (!isRecord(data) || typeof data.hash !== 'string' || typeof data.url !== 'string') {
    throw new TuyileHttpError(500, '分享失败，稍后再试～', null)
  }
  return { hash: data.hash, url: data.url }
}

/** 读取分享（SSR 可用） */
export function apiFetchTuyileShare(hash: string): Promise<ShareDetail> {
  return $fetch<ShareDetail>(`${TUYILE_API_BASE}/share/${encodeURIComponent(hash)}`)
}
