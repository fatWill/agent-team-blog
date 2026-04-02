import { getPool } from '~/server/utils/db'
import type { RowDataPacket } from 'mysql2/promise'

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
interface MessageRow extends RowDataPacket {
  id: number
  device_id: string
  nickname: string | null
  content: string
  last_modified_date: string | null
  ip: string | null
  created_at: Date
  updated_at: Date
}

/**
 * 获取 UTC+8 的今天日期字符串（YYYY-MM-DD）
 */
export function getTodayDateCST(): string {
  const now = new Date()
  // UTC 时间 + 8 小时
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
  // canEdit：是自己的留言 && 今天没有修改过（last_modified_date 为 null 或不等于今天）
  const canEdit = isOwn && (row.last_modified_date === null || row.last_modified_date !== today)

  return {
    id: row.id,
    deviceId: maskDeviceId(row.device_id),
    nickname: row.nickname || '匿名',
    content: row.content,
    isOwn,
    canEdit,
    createdAt: new Date(row.created_at).toISOString(),
    updatedAt: new Date(row.updated_at).toISOString(),
  }
}

/**
 * 获取所有留言列表
 */
export async function getMessageList(currentDeviceId?: string): Promise<MessageItem[]> {
  const pool = getPool()

  const [rows] = await pool.query<MessageRow[]>(
    'SELECT id, device_id, nickname, content, last_modified_date, ip, created_at, updated_at FROM messages ORDER BY created_at DESC',
  )

  return rows.map((row) => rowToMessage(row, currentDeviceId))
}

/**
 * 根据 ID 获取留言
 */
export async function getMessageById(id: number): Promise<MessageRow | null> {
  const pool = getPool()

  const [rows] = await pool.query<MessageRow[]>(
    'SELECT id, device_id, nickname, content, last_modified_date, ip, created_at, updated_at FROM messages WHERE id = ? LIMIT 1',
    [id],
  )

  return rows.length > 0 ? rows[0] : null
}

/**
 * 根据 device_id 获取留言
 */
export async function getMessageByDeviceId(deviceId: string): Promise<MessageRow | null> {
  const pool = getPool()

  const [rows] = await pool.query<MessageRow[]>(
    'SELECT id, device_id, nickname, content, last_modified_date, ip, created_at, updated_at FROM messages WHERE device_id = ? LIMIT 1',
    [deviceId],
  )

  return rows.length > 0 ? rows[0] : null
}

/**
 * 创建留言
 */
export async function createMessage(data: {
  deviceId: string
  nickname?: string
  content: string
  ip?: string
}): Promise<MessageItem> {
  const pool = getPool()

  const [result] = await pool.execute(
    'INSERT INTO messages (device_id, nickname, content, ip) VALUES (?, ?, ?, ?)',
    [data.deviceId, data.nickname || null, data.content, data.ip || null],
  ) as any

  const row = await getMessageById(result.insertId)
  return rowToMessage(row!, data.deviceId)
}

/**
 * 更新留言
 */
export async function updateMessage(id: number, data: {
  nickname?: string
  content: string
  deviceId: string
}): Promise<MessageItem> {
  const pool = getPool()
  const today = getTodayDateCST()

  await pool.execute(
    'UPDATE messages SET nickname = ?, content = ?, last_modified_date = ? WHERE id = ?',
    [data.nickname || null, data.content, today, id],
  )

  const row = await getMessageById(id)
  return rowToMessage(row!, data.deviceId)
}
