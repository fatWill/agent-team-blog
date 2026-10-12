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
