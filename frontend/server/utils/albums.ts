import { getDb } from './db'

/** 相册集数据结构 */
export interface AlbumItem {
  id: number
  name: string
  description: string | null
  coverUrl: string | null
  photoCount: number
  hasPassword: boolean
  createdAt: string
  updatedAt: string
}

/** 照片数据结构 */
export interface PhotoItem {
  id: number
  albumId: number
  url: string
  caption: string | null
  hasPassword: boolean
  createdAt: string
  updatedAt: string
}

/** 数据库行映射 - 相册集 */
interface AlbumRow {
  id: number
  name: string
  description: string | null
  cover_url: string | null
  password_hash: string | null
  photo_count: number
  created_at: string
  updated_at: string
}

/** 数据库行映射 - 照片 */
interface PhotoRow {
  id: number
  album_id: number
  url: string
  caption: string | null
  password_hash: string | null
  created_at: string
  updated_at: string
}

/**
 * 将数据库行转换为 AlbumItem
 */
function rowToAlbum(row: AlbumRow): AlbumItem {
  return {
    id: row.id,
    name: row.name,
    description: row.description ?? null,
    coverUrl: row.cover_url ?? null,
    photoCount: Number(row.photo_count) || 0,
    hasPassword: !!row.password_hash,
    createdAt: row.created_at,
    updatedAt: row.updated_at,
  }
}

/**
 * 将数据库行转换为 PhotoItem
 */
function rowToPhoto(row: PhotoRow): PhotoItem {
  return {
    id: row.id,
    albumId: row.album_id,
    url: row.url,
    caption: row.caption ?? null,
    hasPassword: !!row.password_hash,
    createdAt: row.created_at,
    updatedAt: row.updated_at,
  }
}

/**
 * 获取所有相册集列表，每个相册附带 photo_count 和 cover_url
 */
export function getAlbumList(): AlbumItem[] {
  const db = getDb()

  const sql = `
    SELECT
      a.id, a.name, a.description, a.cover_url, a.password_hash,
      a.created_at, a.updated_at,
      IFNULL(pc.cnt, 0) AS photo_count
    FROM albums a
    LEFT JOIN (
      SELECT album_id, COUNT(*) AS cnt
      FROM photos
      GROUP BY album_id
    ) pc ON pc.album_id = a.id
    ORDER BY a.created_at DESC
  `

  const rows = db.prepare(sql).all() as AlbumRow[]
  return rows.map(rowToAlbum)
}

/**
 * 获取单个相册集详情（含照片数量）
 */
export function getAlbumById(id: number): AlbumItem | null {
  const db = getDb()

  const sql = `
    SELECT
      a.id, a.name, a.description, a.cover_url, a.password_hash,
      a.created_at, a.updated_at,
      IFNULL(pc.cnt, 0) AS photo_count
    FROM albums a
    LEFT JOIN (
      SELECT album_id, COUNT(*) AS cnt
      FROM photos
      WHERE album_id = ?
      GROUP BY album_id
    ) pc ON pc.album_id = a.id
    WHERE a.id = ?
  `

  const row = db.prepare(sql).get(id, id) as AlbumRow | undefined
  if (!row) return null
  return rowToAlbum(row)
}

/**
 * 获取相册下所有照片，按 created_at DESC 排序
 */
export function getPhotosByAlbumId(albumId: number): PhotoItem[] {
  const db = getDb()

  const rows = db.prepare(
    'SELECT id, album_id, url, caption, password_hash, created_at, updated_at FROM photos WHERE album_id = ? ORDER BY created_at DESC',
  ).all(albumId) as PhotoRow[]

  return rows.map(rowToPhoto)
}

/**
 * 创建相册集
 */
export function createAlbum(name: string, description?: string, passwordHash?: string | null): AlbumItem {
  const db = getDb()

  const result = db.prepare(
    'INSERT INTO albums (name, description, password_hash) VALUES (?, ?, ?)',
  ).run(name, description ?? null, passwordHash ?? null)

  return getAlbumById(Number(result.lastInsertRowid))!
}

/**
 * 更新相册集
 */
export function updateAlbum(id: number, data: { name?: string; description?: string; passwordHash?: string | null }): AlbumItem | null {
  const db = getDb()

  const setClauses: string[] = []
  const values: any[] = []

  if (data.name !== undefined) {
    setClauses.push('name = ?')
    values.push(data.name)
  }
  if (data.description !== undefined) {
    setClauses.push('description = ?')
    values.push(data.description)
  }
  if (data.passwordHash !== undefined) {
    setClauses.push('password_hash = ?')
    values.push(data.passwordHash)
  }

  if (setClauses.length === 0) return getAlbumById(id)

  values.push(id)

  const result = db.prepare(
    `UPDATE albums SET ${setClauses.join(', ')} WHERE id = ?`,
  ).run(...values)

  if (result.changes === 0) return null
  return getAlbumById(id)
}

/**
 * 删除相册集（级联删除该相册下所有照片）
 */
export function deleteAlbum(id: number): boolean {
  const db = getDb()

  // 先删除关联照片
  db.prepare('DELETE FROM photos WHERE album_id = ?').run(id)

  // 再删除相册
  const result = db.prepare('DELETE FROM albums WHERE id = ?').run(id)

  return result.changes > 0
}

/**
 * 添加照片到相册
 */
export function addPhoto(albumId: number, url: string, caption?: string, passwordHash?: string | null): PhotoItem {
  const db = getDb()

  const result = db.prepare(
    'INSERT INTO photos (album_id, url, caption, password_hash) VALUES (?, ?, ?, ?)',
  ).run(albumId, url, caption ?? null, passwordHash ?? null)

  const row = db.prepare(
    'SELECT id, album_id, url, caption, password_hash, created_at, updated_at FROM photos WHERE id = ?',
  ).get(Number(result.lastInsertRowid)) as PhotoRow

  return rowToPhoto(row)
}

/**
 * 删除单张照片，返回被删除照片的 album_id（用于后续更新封面）
 */
export function deletePhoto(id: number): { deleted: boolean; albumId: number | null } {
  const db = getDb()

  // 先查询该照片的 album_id
  const photoRow = db.prepare('SELECT album_id FROM photos WHERE id = ?').get(id) as { album_id: number } | undefined

  if (!photoRow) {
    return { deleted: false, albumId: null }
  }

  const albumId = photoRow.album_id

  const result = db.prepare('DELETE FROM photos WHERE id = ?').run(id)

  return { deleted: result.changes > 0, albumId }
}

/**
 * 更新相册封面（取该相册最新照片的 url，无照片则设为 null）
 */
export function updateAlbumCover(albumId: number): void {
  const db = getDb()

  const row = db.prepare(
    'SELECT url FROM photos WHERE album_id = ? ORDER BY created_at DESC LIMIT 1',
  ).get(albumId) as { url: string } | undefined

  const coverUrl = row ? row.url : null

  db.prepare('UPDATE albums SET cover_url = ? WHERE id = ?').run(coverUrl, albumId)
}

/**
 * 获取相册的 password_hash（用于密码验证）
 */
export function getAlbumPasswordHash(id: number): string | null {
  const db = getDb()

  const row = db.prepare('SELECT password_hash FROM albums WHERE id = ? LIMIT 1').get(id) as { password_hash: string | null } | undefined

  if (!row) return null
  return row.password_hash
}

/**
 * 获取照片的 password_hash（用于密码验证）
 */
export function getPhotoPasswordHash(id: number): string | null {
  const db = getDb()

  const row = db.prepare('SELECT password_hash FROM photos WHERE id = ? LIMIT 1').get(id) as { password_hash: string | null } | undefined

  if (!row) return null
  return row.password_hash
}

/**
 * 更新单张照片信息
 */
export function updatePhoto(id: number, data: { caption?: string; passwordHash?: string | null }): PhotoItem | null {
  const db = getDb()

  const setClauses: string[] = []
  const values: any[] = []

  if (data.caption !== undefined) {
    setClauses.push('caption = ?')
    values.push(data.caption)
  }
  if (data.passwordHash !== undefined) {
    setClauses.push('password_hash = ?')
    values.push(data.passwordHash)
  }

  if (setClauses.length === 0) return null

  values.push(id)

  const result = db.prepare(
    `UPDATE photos SET ${setClauses.join(', ')} WHERE id = ?`,
  ).run(...values)

  if (result.changes === 0) return null

  const row = db.prepare(
    'SELECT id, album_id, url, caption, password_hash, created_at, updated_at FROM photos WHERE id = ?',
  ).get(id) as PhotoRow | undefined

  if (!row) return null
  return rowToPhoto(row)
}
