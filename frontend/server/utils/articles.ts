import { getPool } from './db'
import type { RowDataPacket, ResultSetHeader } from 'mysql2/promise'

/** 文章数据结构 */
export interface Article {
  id: string
  title: string
  summary: string
  coverImage: string
  content: Record<string, any>  // Tiptap JSON DSL
  createdAt: string
  updatedAt: string
}

/** 数据库行映射接口 */
interface ArticleRow extends RowDataPacket {
  id: string
  title: string
  summary: string
  cover_image: string
  content: Record<string, any> | string
  created_at: Date
  updated_at: Date
}

/**
 * 将数据库行转换为 Article 对象
 */
function rowToArticle(row: ArticleRow): Article {
  // MySQL JSON 类型会被 mysql2 自动解析为对象，但 fallback 处理字符串情况
  const content = typeof row.content === 'string' ? JSON.parse(row.content) : row.content
  return {
    id: row.id,
    title: row.title,
    summary: row.summary,
    coverImage: row.cover_image || '',
    content,
    createdAt: new Date(row.created_at).toISOString(),
    updatedAt: new Date(row.updated_at).toISOString(),
  }
}

/**
 * 查询文章列表（不含 content），按创建时间倒序
 * @param title 可选，按标题关键词模糊搜索
 */
export async function getArticleList(title?: string): Promise<Omit<Article, 'content'>[]> {
  const pool = getPool()

  let sql = 'SELECT id, title, summary, cover_image, created_at, updated_at FROM articles'
  const params: any[] = []

  if (title && title.trim() !== '') {
    sql += ' WHERE title LIKE ?'
    params.push(`%${title.trim()}%`)
  }

  sql += ' ORDER BY created_at DESC'

  const [rows] = await pool.query<ArticleRow[]>(sql, params)
  return rows.map((row) => ({
    id: row.id,
    title: row.title,
    summary: row.summary,
    coverImage: row.cover_image || '',
    createdAt: new Date(row.created_at).toISOString(),
    updatedAt: new Date(row.updated_at).toISOString(),
  }))
}

/**
 * 根据 ID 查询文章详情
 */
export async function getArticleById(id: string): Promise<Article | null> {
  const pool = getPool()
  const [rows] = await pool.query<ArticleRow[]>(
    'SELECT id, title, summary, cover_image, content, created_at, updated_at FROM articles WHERE id = ?',
    [id],
  )
  if (rows.length === 0) return null
  return rowToArticle(rows[0])
}

/**
 * 创建文章
 */
export async function createArticle(article: {
  id: string
  title: string
  summary: string
  coverImage: string
  content: Record<string, any>
}): Promise<Article> {
  const pool = getPool()
  const now = new Date()
  await pool.query<ResultSetHeader>(
    'INSERT INTO articles (id, title, summary, cover_image, content, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)',
    [article.id, article.title, article.summary, article.coverImage, JSON.stringify(article.content), now, now],
  )
  return {
    id: article.id,
    title: article.title,
    summary: article.summary,
    coverImage: article.coverImage,
    content: article.content,
    createdAt: now.toISOString(),
    updatedAt: now.toISOString(),
  }
}

/**
 * 更新文章（部分更新）
 */
export async function updateArticle(
  id: string,
  updates: { title?: string; summary?: string; coverImage?: string; content?: Record<string, any> },
): Promise<Article | null> {
  const pool = getPool()

  // 构建动态 SET 子句
  const setClauses: string[] = []
  const values: any[] = []

  if (updates.title !== undefined) {
    setClauses.push('title = ?')
    values.push(updates.title)
  }
  if (updates.summary !== undefined) {
    setClauses.push('summary = ?')
    values.push(updates.summary)
  }
  if (updates.coverImage !== undefined) {
    setClauses.push('cover_image = ?')
    values.push(updates.coverImage)
  }
  if (updates.content !== undefined) {
    setClauses.push('content = ?')
    values.push(JSON.stringify(updates.content))
  }

  if (setClauses.length === 0) return null

  // 显式设置 updated_at 更可控
  setClauses.push('updated_at = ?')
  values.push(new Date())
  values.push(id)

  const [result] = await pool.query<ResultSetHeader>(
    `UPDATE articles SET ${setClauses.join(', ')} WHERE id = ?`,
    values,
  )

  if (result.affectedRows === 0) return null

  // 返回更新后的完整文章
  return getArticleById(id)
}

/**
 * 删除文章
 */
export async function deleteArticle(id: string): Promise<boolean> {
  const pool = getPool()
  const [result] = await pool.query<ResultSetHeader>(
    'DELETE FROM articles WHERE id = ?',
    [id],
  )
  return result.affectedRows > 0
}
