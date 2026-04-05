import Database from 'better-sqlite3'

/** 全局数据库实例（惰性初始化，进程级单例） */
let _db: Database.Database | null = null

/**
 * 获取 SQLite 数据库实例（单例）
 * 使用 Nuxt runtimeConfig 读取数据库路径
 */
export function getDb(): Database.Database {
  if (!_db) {
    const config = useRuntimeConfig()
    const dbPath = config.dbPath || process.env.DB_PATH || '/root/blog-data/blog.db'
    _db = new Database(dbPath)
    // WAL 模式提升并发性能
    _db.pragma('journal_mode = WAL')
    _db.pragma('foreign_keys = ON')
  }
  return _db
}
