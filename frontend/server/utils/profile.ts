import type { RowDataPacket, ResultSetHeader } from 'mysql2/promise'
import { getPool } from './db'

/** Profile 数据结构 */
export interface ProfileRow {
  avatar: string
  bio: string
}

interface ProfileDBRow extends RowDataPacket {
  avatar: string
  bio: string
}

/**
 * 获取博主个人资料
 * profile 表只有一条记录（id=1）
 */
export async function getProfile(): Promise<ProfileRow> {
  const pool = getPool()

  const [rows] = await pool.query<ProfileDBRow[]>(
    'SELECT avatar, bio FROM profile WHERE id = 1',
  )

  if (rows.length === 0) {
    // 理论上不会发生，因为初始化时已插入默认记录
    return { avatar: '', bio: '' }
  }

  return {
    avatar: rows[0].avatar,
    bio: rows[0].bio,
  }
}

/**
 * 更新博主个人资料
 * 仅更新传入的字段
 */
export async function updateProfile(data: Partial<ProfileRow>): Promise<ProfileRow> {
  const pool = getPool()

  // 构建动态 SET 子句
  const setClauses: string[] = []
  const values: (string)[] = []

  if (data.avatar !== undefined) {
    setClauses.push('avatar = ?')
    values.push(data.avatar)
  }

  if (data.bio !== undefined) {
    setClauses.push('bio = ?')
    values.push(data.bio)
  }

  if (setClauses.length > 0) {
    await pool.query<ResultSetHeader>(
      `UPDATE profile SET ${setClauses.join(', ')} WHERE id = 1`,
      values,
    )
  }

  // 返回更新后的完整数据
  return getProfile()
}
