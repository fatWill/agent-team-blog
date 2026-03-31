import { getPool } from './db'
import type { RowDataPacket } from 'mysql2/promise'

/** 更新日志数据结构 */
export interface ChangelogItem {
  version: string
  date: string
  logs: string[]
  createdAt: string
  updatedAt: string
}

/** 数据库行映射接口 */
interface ChangelogRow extends RowDataPacket {
  id: number
  version: string
  date: string
  logs: string[] | string
  created_at: Date
  updated_at: Date
}

/**
 * 将数据库行转换为 ChangelogItem
 */
function rowToChangelog(row: ChangelogRow): ChangelogItem {
  const logs = typeof row.logs === 'string' ? JSON.parse(row.logs) : row.logs
  return {
    version: row.version,
    date: row.date,
    logs,
    createdAt: new Date(row.created_at).toISOString(),
    updatedAt: new Date(row.updated_at).toISOString(),
  }
}

/**
 * 查询所有更新日志，按版本号倒序（id 倒序即为时间倒序）
 */
export async function getChangelogList(): Promise<ChangelogItem[]> {
  const pool = getPool()
  const [rows] = await pool.query<ChangelogRow[]>(
    'SELECT id, version, date, logs, created_at, updated_at FROM changelogs ORDER BY id DESC',
  )
  return rows.map(rowToChangelog)
}
