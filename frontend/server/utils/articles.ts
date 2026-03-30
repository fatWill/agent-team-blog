import { promises as fs } from 'node:fs'
import { resolve, dirname } from 'node:path'

/** 文章数据结构 */
export interface Article {
  id: string
  title: string
  summary: string
  content: Record<string, any>  // Tiptap JSON DSL
  createdAt: string
  updatedAt: string
}

/** JSON 文件存储结构 */
interface ArticlesData {
  articles: Article[]
}

/** 数据文件路径 */
const DATA_FILE = resolve(process.cwd(), 'server/data/articles.json')

/**
 * 确保数据目录和文件存在，不存在则自动创建
 */
async function ensureDataFile(): Promise<void> {
  try {
    await fs.access(DATA_FILE)
  } catch {
    // 文件不存在，创建目录和初始文件
    await fs.mkdir(dirname(DATA_FILE), { recursive: true })
    await fs.writeFile(DATA_FILE, JSON.stringify({ articles: [] }, null, 2), 'utf-8')
  }
}

/**
 * 读取所有文章数据
 */
export async function readArticles(): Promise<Article[]> {
  await ensureDataFile()
  const raw = await fs.readFile(DATA_FILE, 'utf-8')
  const data: ArticlesData = JSON.parse(raw)
  return data.articles
}

/**
 * 写入文章数据到 JSON 文件
 */
export async function writeArticles(articles: Article[]): Promise<void> {
  await ensureDataFile()
  const data: ArticlesData = { articles }
  await fs.writeFile(DATA_FILE, JSON.stringify(data, null, 2), 'utf-8')
}
