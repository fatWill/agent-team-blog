import { getDb } from './db'

/** 更新日志数据结构 */
export interface ChangelogItem {
  version: string
  date: string
  logs: string[]
  createdAt: string
  updatedAt: string
}

/** 数据库行映射接口 */
interface ChangelogRow {
  id: number
  version: string
  date: string
  logs: string
  created_at: string
  updated_at: string
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
    createdAt: row.created_at,
    updatedAt: row.updated_at,
  }
}

/**
 * 查询所有更新日志，按版本号倒序（id 倒序即为时间倒序）
 */
export function getChangelogList(): ChangelogItem[] {
  const db = getDb()
  const rows = db.prepare(
    'SELECT id, version, date, logs, created_at, updated_at FROM changelogs ORDER BY id DESC',
  ).all() as ChangelogRow[]
  return rows.map(rowToChangelog)
}
