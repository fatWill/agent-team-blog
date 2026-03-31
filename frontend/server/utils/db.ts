import mysql from 'mysql2/promise'
import type { Pool } from 'mysql2/promise'

/** 全局连接池（惰性初始化，进程级单例） */
let pool: Pool | null = null

/**
 * 获取 MySQL 连接池（单例）
 * 使用 Nuxt runtimeConfig 读取数据库配置
 */
export function getPool(): Pool {
  if (!pool) {
    const config = useRuntimeConfig()
    pool = mysql.createPool({
      host: config.dbHost,
      port: Number(config.dbPort) || 3306,
      user: config.dbUser,
      password: config.dbPassword,
      database: config.dbName,
      // 2核2G 服务器，连接数保守设置
      connectionLimit: 5,
      // 等待连接超时 10 秒
      waitForConnections: true,
      queueLimit: 0,
      // 空闲连接回收
      idleTimeout: 60000,
      // 启用 JSON 自动解析
      typeCast: true,
    })
  }
  return pool
}
