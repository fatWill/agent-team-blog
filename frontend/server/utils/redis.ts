import Redis from 'ioredis'

let client: Redis | null = null

/**
 * 获取 Redis 客户端单例
 */
export function getRedis(): Redis {
  if (!client) {
    const config = useRuntimeConfig()
    client = new Redis({
      host: config.redisHost || '127.0.0.1',
      port: Number(config.redisPort) || 6379,
      password: config.redisPassword || undefined,
      lazyConnect: true,
      // 连接失败不抛出，静默降级
      retryStrategy: (times) => Math.min(times * 100, 3000),
    })
    client.on('error', (err) => {
      console.error('[Redis] connection error:', err.message)
    })
  }
  return client
}
