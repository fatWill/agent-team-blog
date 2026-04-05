import { getDb } from './db'

/** 文章数据结构 */
export interface Article {
  id: string
  title: string
  summary: string
  coverImage: string
  content: Record<string, any>  // Tiptap JSON DSL
  likeCount: number
  createdAt: string
  updatedAt: string
}

/** 数据库行映射接口 */
interface ArticleRow {
  id: string
  title: string
  summary: string
  cover_image: string
  content: string
  like_count: number
  created_at: string
  updated_at: string
}

/**
 * 将数据库行转换为 Article 对象
 */
function rowToArticle(row: ArticleRow): Article {
  const content = typeof row.content === 'string' ? JSON.parse(row.content) : row.content
  return {
    id: row.id,
    title: row.title,
    summary: row.summary,
    coverImage: row.cover_image || '',
    content,
    likeCount: Number(row.like_count) || 0,
    createdAt: row.created_at,
    updatedAt: row.updated_at,
  }
}

/**
 * 查询文章列表（不含 content），按创建时间倒序
 * @param title 可选，按标题关键词模糊搜索
 */
export function getArticleList(title?: string): Omit<Article, 'content'>[] {
  const db = getDb()

  let sql = 'SELECT id, title, summary, cover_image, like_count, created_at, updated_at FROM articles'
  const params: any[] = []

  if (title && title.trim() !== '') {
    sql += ' WHERE title LIKE ?'
    params.push(`%${title.trim()}%`)
  }

  sql += ' ORDER BY created_at DESC'

  const rows = db.prepare(sql).all(...params) as ArticleRow[]
  return rows.map((row) => ({
    id: row.id,
    title: row.title,
    summary: row.summary,
    coverImage: row.cover_image || '',
    likeCount: Number(row.like_count) || 0,
    createdAt: row.created_at,
    updatedAt: row.updated_at,
  }))
}

/**
 * 根据 ID 查询文章详情
 */
export function getArticleById(id: string): Article | null {
  const db = getDb()
  const row = db.prepare(
    'SELECT id, title, summary, cover_image, content, like_count, created_at, updated_at FROM articles WHERE id = ?',
  ).get(id) as ArticleRow | undefined

  if (!row) return null
  return rowToArticle(row)
}

/**
 * 创建文章
 */
export function createArticle(article: {
  id: string
  title: string
  summary: string
  coverImage: string
  content: Record<string, any>
}): Article {
  const db = getDb()
  const now = new Date().toISOString()
  db.prepare(
    'INSERT INTO articles (id, title, summary, cover_image, content, like_count, created_at, updated_at) VALUES (?, ?, ?, ?, ?, 0, ?, ?)',
  ).run(article.id, article.title, article.summary, article.coverImage, JSON.stringify(article.content), now, now)

  return {
    id: article.id,
    title: article.title,
    summary: article.summary,
    coverImage: article.coverImage,
    content: article.content,
    likeCount: 0,
    createdAt: now,
    updatedAt: now,
  }
}

/**
 * 更新文章（部分更新）
 */
export function updateArticle(
  id: string,
  updates: { title?: string; summary?: string; coverImage?: string; content?: Record<string, any> },
): Article | null {
  const db = getDb()

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

  setClauses.push('updated_at = ?')
  values.push(new Date().toISOString())
  values.push(id)

  const result = db.prepare(
    `UPDATE articles SET ${setClauses.join(', ')} WHERE id = ?`,
  ).run(...values)

  if (result.changes === 0) return null

  return getArticleById(id)
}

/**
 * 删除文章
 */
export function deleteArticle(id: string): boolean {
  const db = getDb()
  const result = db.prepare('DELETE FROM articles WHERE id = ?').run(id)
  return result.changes > 0
}
