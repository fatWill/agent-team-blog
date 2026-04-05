import { getDb } from '~/server/utils/db'

/** 留言数据结构 */
export interface MessageItem {
  id: number
  deviceId: string
  nickname: string
  content: string
  isOwn: boolean
  canEdit: boolean
  createdAt: string
  updatedAt: string
}

/** 数据库行映射 */
interface MessageRow {
  id: number
  device_id: string
  nickname: string | null
  content: string
  last_modified_date: string | null
  ip: string | null
  created_at: string
  updated_at: string
}

/**
 * 获取 UTC+8 的今天日期字符串（YYYY-MM-DD）
 */
export function getTodayDateCST(): string {
  const now = new Date()
  const cst = new Date(now.getTime() + 8 * 60 * 60 * 1000)
  return cst.toISOString().slice(0, 10)
}

/**
 * 脱敏 deviceId：只显示前 4 位 + ****
 */
function maskDeviceId(deviceId: string): string {
  if (deviceId.length <= 4) return '****'
  return deviceId.slice(0, 4) + '****'
}

/**
 * 将数据库行转换为 MessageItem
 * @param row 数据库行
 * @param currentDeviceId 当前请求的 deviceId（用于判断 isOwn 和 canEdit）
 */
function rowToMessage(row: MessageRow, currentDeviceId?: string): MessageItem {
  const isOwn = !!currentDeviceId && row.device_id === currentDeviceId
  const today = getTodayDateCST()
  const canEdit = isOwn && (row.last_modified_date === null || row.last_modified_date !== today)

  return {
    id: row.id,
    deviceId: maskDeviceId(row.device_id),
    nickname: row.nickname || '匿名',
    content: row.content,
    isOwn,
    canEdit,
    createdAt: row.created_at,
    updatedAt: row.updated_at,
  }
}

/**
 * 获取所有留言列表
 */
export function getMessageList(currentDeviceId?: string): MessageItem[] {
  const db = getDb()

  const rows = db.prepare(
    'SELECT id, device_id, nickname, content, last_modified_date, ip, created_at, updated_at FROM messages ORDER BY created_at DESC',
  ).all() as MessageRow[]

  return rows.map((row) => rowToMessage(row, currentDeviceId))
}

/**
 * 根据 ID 获取留言
 */
export function getMessageById(id: number): MessageRow | null {
  const db = getDb()

  const row = db.prepare(
    'SELECT id, device_id, nickname, content, last_modified_date, ip, created_at, updated_at FROM messages WHERE id = ? LIMIT 1',
  ).get(id) as MessageRow | undefined

  return row ?? null
}

/**
 * 根据 device_id 获取留言
 */
export function getMessageByDeviceId(deviceId: string): MessageRow | null {
  const db = getDb()

  const row = db.prepare(
    'SELECT id, device_id, nickname, content, last_modified_date, ip, created_at, updated_at FROM messages WHERE device_id = ? LIMIT 1',
  ).get(deviceId) as MessageRow | undefined

  return row ?? null
}

/**
 * 创建留言
 */
export function createMessage(data: {
  deviceId: string
  nickname?: string
  content: string
  ip?: string
}): MessageItem {
  const db = getDb()

  const result = db.prepare(
    'INSERT INTO messages (device_id, nickname, content, ip) VALUES (?, ?, ?, ?)',
  ).run(data.deviceId, data.nickname || null, data.content, data.ip || null)

  const row = getMessageById(Number(result.lastInsertRowid))
  return rowToMessage(row!, data.deviceId)
}

/**
 * 更新留言
 */
export function updateMessage(id: number, data: {
  nickname?: string
  content: string
  deviceId: string
}): MessageItem {
  const db = getDb()
  const today = getTodayDateCST()

  db.prepare(
    'UPDATE messages SET nickname = ?, content = ?, last_modified_date = ? WHERE id = ?',
  ).run(data.nickname || null, data.content, today, id)

  const row = getMessageById(id)
  return rowToMessage(row!, data.deviceId)
}
