import { getPool } from './db'
import type { RowDataPacket, ResultSetHeader } from 'mysql2/promise'

/** 文章数据结构 */
export interface Article {
  id: string
  title: string
  summary: string
  content: Record<string, any>  // Tiptap JSON DSL
  createdAt: string
  updatedAt: string
}

/** 数据库行映射接口 */
interface ArticleRow extends RowDataPacket {
  id: string
  title: string
  summary: string
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
    content,
    createdAt: new Date(row.created_at).toISOString(),
    updatedAt: new Date(row.updated_at).toISOString(),
  }
}

/**
 * 查询所有文章（不含 content），按创建时间倒序
 */
export async function getArticleList(): Promise<Omit<Article, 'content'>[]> {
  const pool = getPool()
  const [rows] = await pool.query<ArticleRow[]>(
    'SELECT id, title, summary, created_at, updated_at FROM articles ORDER BY created_at DESC',
  )
  return rows.map((row) => ({
    id: row.id,
    title: row.title,
    summary: row.summary,
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
    'SELECT id, title, summary, content, created_at, updated_at FROM articles WHERE id = ?',
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
  content: Record<string, any>
}): Promise<Article> {
  const pool = getPool()
  const now = new Date()
  await pool.query<ResultSetHeader>(
    'INSERT INTO articles (id, title, summary, content, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)',
    [article.id, article.title, article.summary, JSON.stringify(article.content), now, now],
  )
  return {
    id: article.id,
    title: article.title,
    summary: article.summary,
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
  updates: { title?: string; summary?: string; content?: Record<string, any> },
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
  if (updates.content !== undefined) {
    setClauses.push('content = ?')
    values.push(JSON.stringify(updates.content))
  }

  if (setClauses.length === 0) return null

  // updated_at 由 MySQL ON UPDATE CURRENT_TIMESTAMP 自动处理，但显式设置更可控
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
