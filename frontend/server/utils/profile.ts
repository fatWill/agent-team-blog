import { getDb } from './db'

/** Profile 数据结构 */
export interface ProfileRow {
  avatar: string
  bio: string
}

/**
 * 获取博主个人资料
 * profile 表只有一条记录（id=1）
 */
export function getProfile(): ProfileRow {
  const db = getDb()

  const row = db.prepare('SELECT avatar, bio FROM profile WHERE id = 1').get() as { avatar: string; bio: string } | undefined

  if (!row) {
    return { avatar: '', bio: '' }
  }

  return {
    avatar: row.avatar,
    bio: row.bio,
  }
}

/**
 * 更新博主个人资料
 * 仅更新传入的字段
 */
export function updateProfile(data: Partial<ProfileRow>): ProfileRow {
  const db = getDb()

  const setClauses: string[] = []
  const values: string[] = []

  if (data.avatar !== undefined) {
    setClauses.push('avatar = ?')
    values.push(data.avatar)
  }

  if (data.bio !== undefined) {
    setClauses.push('bio = ?')
    values.push(data.bio)
  }

  if (setClauses.length > 0) {
    db.prepare(`UPDATE profile SET ${setClauses.join(', ')} WHERE id = 1`).run(...values)
  }

  return getProfile()
}
