import type { AgUiEvent, AgUiMessage, ChatBubble, ChatStatus, TuyileStatusValue } from '../types'
import { apiStreamTuyileChat, TuyileHttpError, TUYILE_MAX_MESSAGE_LENGTH } from '../api'
import { genId } from '../utils'

function isStatusValue(v: unknown): v is TuyileStatusValue {
  if (typeof v !== 'object' || v === null) return false
  const r = v as Record<string, unknown>
  return typeof r.text === 'string' && r.text !== '' && typeof r.stage === 'string'
    && (r.range === undefined || typeof r.range === 'string')
}

/** 请求体安全上限（后端 64KB），超出时从最早的消息开始裁剪 */
const MAX_BODY_BYTES = 60 * 1024

export interface ChatNotice {
  kind: 'error' | 'rate'
  message: string
  /** 429/503 时的剩余等待秒数 */
  retryAfter: number | null
}

/**
 * 图一乐 AG-UI 聊天状态。
 * 对话只存在于组件内存，刷新/离开页面即清空，不落 localStorage / sessionStorage（不持久化姓名生日）。
 */
export function useTuyileChat() {
  const threadId = genId('thr')
  const bubbles = ref<ChatBubble[]>([])
  const isStreaming = ref(false)
  const notice = ref<ChatNotice | null>(null)
  const cooldown = ref(0)
  const started = ref(false)
  const status = ref<ChatStatus | null>(null)

  let controller: AbortController | null = null
  let cooldownTimer: ReturnType<typeof setInterval> | null = null

  function startCooldown(seconds: number) {
    if (cooldownTimer) clearInterval(cooldownTimer)
    cooldown.value = seconds
    cooldownTimer = setInterval(() => {
      cooldown.value -= 1
      if (cooldown.value <= 0) {
        cooldown.value = 0
        if (cooldownTimer) clearInterval(cooldownTimer)
        cooldownTimer = null
        if (notice.value?.kind === 'rate') notice.value = null
      }
    }, 1000)
  }

  /** 已完成气泡 → AG-UI messages（丢弃空内容），超出体积上限时裁掉最早的消息 */
  function buildHistory(): AgUiMessage[] {
    const list: AgUiMessage[] = bubbles.value
      .filter(b => b.content.trim() !== '')
      .map(b => ({ id: b.id, role: b.role, content: b.content }))
    const encoder = new TextEncoder()
    while (list.length > 1 && encoder.encode(JSON.stringify(list)).length > MAX_BODY_BYTES) {
      list.shift()
    }
    return list
  }

  function findBubble(id: string): ChatBubble | undefined {
    return bubbles.value.find(b => b.id === id)
  }

  /**
   * 发起一次 run。
   * @returns 失败且未收到任何回复时返回 false（调用方可回滚输入）
   */
  async function run(): Promise<boolean> {
    if (isStreaming.value) return false
    notice.value = null
    isStreaming.value = true
    controller = new AbortController()
    const touched: string[] = []
    let gotContent = false

    const onEvent = (ev: AgUiEvent) => {
      switch (ev.type) {
        case 'RUN_STARTED': {
          if (!status.value) status.value = { stage: 'thinking', text: '掐指一算中…' }
          break
        }
        case 'CUSTOM': {
          if (ev.name === 'tuyile.status' && isStatusValue(ev.value)) {
            status.value = { stage: ev.value.stage, text: ev.value.text, range: ev.value.range }
          }
          break
        }
        case 'REASONING_MESSAGE_CONTENT':
        case 'REASONING_START':
        case 'REASONING_END':
        case 'REASONING_MESSAGE_START':
        case 'REASONING_MESSAGE_END':
          // 状态摘要统一由 CUSTOM tuyile.status 驱动；reasoning 不渲染、不回传
          break
        case 'TEXT_MESSAGE_START': {
          status.value = null
          if (!findBubble(ev.messageId)) {
            bubbles.value.push({ id: ev.messageId, role: 'assistant', content: '', streaming: true })
            touched.push(ev.messageId)
          }
          break
        }
        case 'TEXT_MESSAGE_CONTENT': {
          status.value = null
          let b = findBubble(ev.messageId)
          if (!b) {
            bubbles.value.push({ id: ev.messageId, role: 'assistant', content: '', streaming: true })
            touched.push(ev.messageId)
            b = findBubble(ev.messageId)
          }
          if (b && ev.delta) {
            b.content += ev.delta
            gotContent = true
          }
          break
        }
        case 'TEXT_MESSAGE_END': {
          const b = findBubble(ev.messageId)
          if (b) b.streaming = false
          break
        }
        case 'RUN_ERROR': {
          status.value = null
          notice.value = { kind: 'error', message: ev.message || '图一乐走神了，稍后再试～', retryAfter: null }
          break
        }
        default:
          break
      }
    }

    let ok = true
    try {
      await apiStreamTuyileChat(
        { threadId, runId: genId('run'), messages: buildHistory() },
        onEvent,
        controller.signal,
      )
    } catch (err: unknown) {
      if (err instanceof DOMException && err.name === 'AbortError') {
        // 用户主动停止：保留已生成的部分
      } else if (err instanceof TuyileHttpError) {
        ok = false
        const waitable = err.status === 429 || err.status === 503
        notice.value = {
          kind: waitable ? 'rate' : 'error',
          message: err.message,
          retryAfter: waitable ? err.retryAfter : null,
        }
        if (waitable && err.retryAfter) startCooldown(err.retryAfter)
      } else {
        ok = false
        notice.value = { kind: 'error', message: '网络开小差了，检查一下网络再试～', retryAfter: null }
      }
    } finally {
      status.value = null
      for (const id of touched) {
        const b = findBubble(id)
        if (b) b.streaming = false
      }
      // 清理空的 assistant 气泡
      bubbles.value = bubbles.value.filter(b => !(b.role === 'assistant' && b.content === '' && touched.includes(b.id)))
      isStreaming.value = false
      controller = null
    }
    if (notice.value?.kind === 'error' && !gotContent) ok = false
    return ok && (gotContent || notice.value === null)
  }

  /** 进入聊天：messages: [] 让图一乐主动打招呼 */
  async function greet() {
    if (started.value) return
    started.value = true
    await run()
  }

  /** 重试打招呼（首次失败时） */
  async function retryGreet() {
    if (bubbles.value.length) return
    await run()
  }

  /**
   * 发送用户消息。
   * @returns 失败且需要回滚时返回原文本，供输入框恢复；成功返回 null
   */
  async function send(text: string): Promise<string | null> {
    const content = text.trim()
    if (!content || isStreaming.value || cooldown.value > 0) return null
    if (Array.from(content).length > TUYILE_MAX_MESSAGE_LENGTH) {
      notice.value = { kind: 'error', message: `单条消息最多 ${TUYILE_MAX_MESSAGE_LENGTH} 字`, retryAfter: null }
      return content
    }
    const userId = genId('u')
    bubbles.value.push({ id: userId, role: 'user', content, streaming: false })
    const ok = await run()
    if (!ok) {
      // 未得到任何回复：撤回这条用户消息，避免历史中出现连续 user，并把文本还给输入框
      bubbles.value = bubbles.value.filter(b => b.id !== userId)
      return content
    }
    return null
  }

  function stop() {
    controller?.abort()
  }

  function dismissNotice() {
    if (notice.value?.kind === 'rate' && cooldown.value > 0) return
    notice.value = null
  }

  onBeforeUnmount(() => {
    controller?.abort()
    if (cooldownTimer) clearInterval(cooldownTimer)
  })

  return { bubbles, isStreaming, notice, cooldown, status, greet, retryGreet, send, stop, dismissNotice }
}
