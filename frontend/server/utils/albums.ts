import { getPool } from './db'
import type { RowDataPacket, ResultSetHeader } from 'mysql2/promise'

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
interface AlbumRow extends RowDataPacket {
  id: number
  name: string
  description: string | null
  cover_url: string | null
  password_hash: string | null
  photo_count: number
  created_at: Date
  updated_at: Date
}

/** 数据库行映射 - 照片 */
interface PhotoRow extends RowDataPacket {
  id: number
  album_id: number
  url: string
  caption: string | null
  password_hash: string | null
  created_at: Date
  updated_at: Date
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
    createdAt: new Date(row.created_at).toISOString(),
    updatedAt: new Date(row.updated_at).toISOString(),
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
    createdAt: new Date(row.created_at).toISOString(),
    updatedAt: new Date(row.updated_at).toISOString(),
  }
}

/**
 * 获取所有相册集列表，每个相册附带 photo_count 和 cover_url
 * cover_url 取该相册最新一张照片的 url，无照片时为 null
 */
export async function getAlbumList(): Promise<AlbumItem[]> {
  const pool = getPool()

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

  const [rows] = await pool.query<AlbumRow[]>(sql)
  return rows.map(rowToAlbum)
}

/**
 * 获取单个相册集详情（含照片数量）
 */
export async function getAlbumById(id: number): Promise<AlbumItem | null> {
  const pool = getPool()

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

  const [rows] = await pool.query<AlbumRow[]>(sql, [id, id])
  if (rows.length === 0) return null
  return rowToAlbum(rows[0])
}

/**
 * 获取相册下所有照片，按 created_at DESC 排序
 */
export async function getPhotosByAlbumId(albumId: number): Promise<PhotoItem[]> {
  const pool = getPool()

  const [rows] = await pool.query<PhotoRow[]>(
    'SELECT id, album_id, url, caption, password_hash, created_at, updated_at FROM photos WHERE album_id = ? ORDER BY created_at DESC',
    [albumId],
  )
  return rows.map(rowToPhoto)
}

/**
 * 创建相册集
 */
export async function createAlbum(name: string, description?: string, passwordHash?: string | null): Promise<AlbumItem> {
  const pool = getPool()

  const [result] = await pool.query<ResultSetHeader>(
    'INSERT INTO albums (name, description, password_hash) VALUES (?, ?, ?)',
    [name, description ?? null, passwordHash ?? null],
  )

  return (await getAlbumById(result.insertId))!
}

/**
 * 更新相册集
 */
export async function updateAlbum(id: number, data: { name?: string; description?: string; passwordHash?: string | null }): Promise<AlbumItem | null> {
  const pool = getPool()

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

  const [result] = await pool.query<ResultSetHeader>(
    `UPDATE albums SET ${setClauses.join(', ')} WHERE id = ?`,
    values,
  )

  if (result.affectedRows === 0) return null
  return getAlbumById(id)
}

/**
 * 删除相册集（级联删除该相册下所有照片）
 */
export async function deleteAlbum(id: number): Promise<boolean> {
  const pool = getPool()

  // 先删除关联照片
  await pool.query<ResultSetHeader>(
    'DELETE FROM photos WHERE album_id = ?',
    [id],
  )

  // 再删除相册
  const [result] = await pool.query<ResultSetHeader>(
    'DELETE FROM albums WHERE id = ?',
    [id],
  )

  return result.affectedRows > 0
}

/**
 * 添加照片到相册
 */
export async function addPhoto(albumId: number, url: string, caption?: string, passwordHash?: string | null): Promise<PhotoItem> {
  const pool = getPool()

  const [result] = await pool.query<ResultSetHeader>(
    'INSERT INTO photos (album_id, url, caption, password_hash) VALUES (?, ?, ?, ?)',
    [albumId, url, caption ?? null, passwordHash ?? null],
  )

  const [rows] = await pool.query<PhotoRow[]>(
    'SELECT id, album_id, url, caption, password_hash, created_at, updated_at FROM photos WHERE id = ?',
    [result.insertId],
  )

  return rowToPhoto(rows[0])
}

/**
 * 删除单张照片，返回被删除照片的 album_id（用于后续更新封面）
 */
export async function deletePhoto(id: number): Promise<{ deleted: boolean; albumId: number | null }> {
  const pool = getPool()

  // 先查询该照片的 album_id
  const [photoRows] = await pool.query<PhotoRow[]>(
    'SELECT album_id FROM photos WHERE id = ?',
    [id],
  )

  if (photoRows.length === 0) {
    return { deleted: false, albumId: null }
  }

  const albumId = photoRows[0].album_id

  const [result] = await pool.query<ResultSetHeader>(
    'DELETE FROM photos WHERE id = ?',
    [id],
  )

  return { deleted: result.affectedRows > 0, albumId }
}

/**
 * 更新相册封面（取该相册最新照片的 url，无照片则设为 null）
 */
export async function updateAlbumCover(albumId: number): Promise<void> {
  const pool = getPool()

  // 查询该相册最新一张照片
  const [rows] = await pool.query<PhotoRow[]>(
    'SELECT url FROM photos WHERE album_id = ? ORDER BY created_at DESC LIMIT 1',
    [albumId],
  )

  const coverUrl = rows.length > 0 ? rows[0].url : null

  await pool.query<ResultSetHeader>(
    'UPDATE albums SET cover_url = ? WHERE id = ?',
    [coverUrl, albumId],
  )
}

/**
 * 获取相册的 password_hash（用于密码验证）
 */
export async function getAlbumPasswordHash(id: number): Promise<string | null> {
  const pool = getPool()

  const [rows] = await pool.query<(RowDataPacket & { password_hash: string | null })[]>(
    'SELECT password_hash FROM albums WHERE id = ? LIMIT 1',
    [id],
  )

  if (rows.length === 0) return null
  return rows[0].password_hash
}

/**
 * 获取照片的 password_hash（用于密码验证）
 */
export async function getPhotoPasswordHash(id: number): Promise<string | null> {
  const pool = getPool()

  const [rows] = await pool.query<(RowDataPacket & { password_hash: string | null })[]>(
    'SELECT password_hash FROM photos WHERE id = ? LIMIT 1',
    [id],
  )

  if (rows.length === 0) return null
  return rows[0].password_hash
}

/**
 * 更新单张照片信息
 */
export async function updatePhoto(id: number, data: { caption?: string; passwordHash?: string | null }): Promise<PhotoItem | null> {
  const pool = getPool()

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

  const [result] = await pool.query<ResultSetHeader>(
    `UPDATE photos SET ${setClauses.join(', ')} WHERE id = ?`,
    values,
  )

  if (result.affectedRows === 0) return null

  const [rows] = await pool.query<PhotoRow[]>(
    'SELECT id, album_id, url, caption, password_hash, created_at, updated_at FROM photos WHERE id = ?',
    [id],
  )

  if (rows.length === 0) return null
  return rowToPhoto(rows[0])
}
