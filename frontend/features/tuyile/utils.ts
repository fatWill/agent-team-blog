import MarkdownIt from 'markdown-it'
import { ALMANAC_MAX_DATE, ALMANAC_MIN_DATE } from './api'

const BEIJING_TZ = 'Asia/Shanghai'

/** 北京时间今天 YYYY-MM-DD */
export function beijingToday(): string {
  return new Intl.DateTimeFormat('en-CA', {
    timeZone: BEIJING_TZ,
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
  }).format(new Date())
}

/** 北京时间当前小时 0-23 */
export function beijingHour(): number {
  const h = new Intl.DateTimeFormat('en-GB', { timeZone: BEIJING_TZ, hour: '2-digit', hourCycle: 'h23' }).format(new Date())
  return Number(h) % 24
}

/** 小时 → 时辰下标（子=0 … 亥=11，子时 23:00-00:59） */
export function hourToShiChenIndex(hour: number): number {
  return Math.floor(((hour + 1) % 24) / 2)
}

export function isValidDateStr(s: string): boolean {
  if (!/^\d{4}-\d{2}-\d{2}$/.test(s)) return false
  const [y, m, d] = s.split('-').map(Number)
  const dt = new Date(Date.UTC(y, m - 1, d))
  return dt.getUTCFullYear() === y && dt.getUTCMonth() === m - 1 && dt.getUTCDate() === d
}

function pad(n: number): string {
  return String(n).padStart(2, '0')
}

/** 日期加减天数，并夹在 1900-2100 支持范围内 */
export function addDays(dateStr: string, delta: number): string {
  const [y, m, d] = dateStr.split('-').map(Number)
  const dt = new Date(Date.UTC(y, m - 1, d + delta))
  const out = `${dt.getUTCFullYear()}-${pad(dt.getUTCMonth() + 1)}-${pad(dt.getUTCDate())}`
  if (out < ALMANAC_MIN_DATE) return ALMANAC_MIN_DATE
  if (out > ALMANAC_MAX_DATE) return ALMANAC_MAX_DATE
  return out
}

/** YYYY-MM 加减月份 */
export function addMonths(month: string, delta: number): string {
  const [y, m] = month.split('-').map(Number)
  const dt = new Date(Date.UTC(y, m - 1 + delta, 1))
  return `${dt.getUTCFullYear()}-${pad(dt.getUTCMonth() + 1)}`
}

/** 某月 1 号是周几（周一=0 … 周日=6） */
export function mondayOffset(month: string): number {
  const [y, m] = month.split('-').map(Number)
  return (new Date(Date.UTC(y, m - 1, 1)).getUTCDay() + 6) % 7
}

export function formatCnDate(dateStr: string): string {
  const [y, m, d] = dateStr.split('-').map(Number)
  return `${y}年${m}月${d}日`
}

export function genId(prefix: string): string {
  const c = globalThis.crypto
  if (c && typeof c.randomUUID === 'function') return `${prefix}-${c.randomUUID()}`
  return `${prefix}-${Date.now().toString(36)}-${Math.random().toString(36).slice(2, 10)}`
}

// ====== Markdown（禁用 HTML，防 XSS） ======
const md = new MarkdownIt({ html: false, linkify: true, breaks: true })
const defaultLinkOpen = md.renderer.rules.link_open
  || ((tokens, idx, options, _env, self) => self.renderToken(tokens, idx, options))
md.renderer.rules.link_open = (tokens, idx, options, env, self) => {
  tokens[idx].attrSet('target', '_blank')
  tokens[idx].attrSet('rel', 'noopener noreferrer nofollow')
  return defaultLinkOpen(tokens, idx, options, env, self)
}

export function renderMarkdown(src: string): string {
  return md.render(src)
}

// ====== 分享内容裁剪 ======

export interface SharePlan {
  messages: { role: 'user' | 'assistant', content: string }[]
  total: number
  /** 因条数/体积上限从开头丢弃的条数 */
  dropped: number
  /** 被截断到 4000 字的条数 */
  truncated: number
}

/**
 * 按后端限制裁剪待分享消息：只保留 user/assistant 非空消息，单条截断 maxLen 字，
 * 最多保留最近 maxCount 条，总体积超过 maxBytes 时继续丢最早的；保证首条之后仍含 assistant。
 */
export function planShare(
  source: { role: string, content: string }[],
  maxCount: number,
  maxLen: number,
  maxBytes: number,
): SharePlan {
  let truncated = 0
  const clean = source
    .filter((m): m is { role: 'user' | 'assistant', content: string } =>
      (m.role === 'user' || m.role === 'assistant') && m.content.trim() !== '')
    .map((m) => {
      const chars = Array.from(m.content)
      if (chars.length > maxLen) {
        truncated++
        return { role: m.role, content: chars.slice(0, maxLen).join('') }
      }
      return { role: m.role, content: m.content }
    })
  let list = clean.slice(-maxCount)
  const encoder = new TextEncoder()
  while (list.length > 1 && encoder.encode(JSON.stringify({ messages: list })).length > maxBytes) {
    list = list.slice(1)
  }
  return { messages: list, total: clean.length, dropped: clean.length - list.length, truncated }
}

/** Markdown 转纯文本摘要（用于 SEO / 分享卡片描述） */
export function markdownSummary(src: string, max = 80): string {
  const text = src
    .replace(/```[\s\S]*?```/g, ' ')
    .replace(/!\[[^\]]*\]\([^)]*\)/g, ' ')
    .replace(/\[([^\]]*)\]\([^)]*\)/g, '$1')
    .replace(/^\s*[-*_]{3,}\s*$/gm, ' ')
    .replace(/^\s{0,3}(?:#{1,6}|>+|[-*+]|\d+\.)\s+/gm, '')
    .replace(/(\*\*|__|~~|`)/g, '')
    .replace(/\s+/g, ' ')
    .trim()
  const chars = Array.from(text)
  return chars.length > max ? `${chars.slice(0, max).join('')}…` : text
}

/** 是否微信内置浏览器 */
export function isWechatUA(): boolean {
  if (import.meta.server) return false
  return /MicroMessenger/i.test(navigator.userAgent)
}
