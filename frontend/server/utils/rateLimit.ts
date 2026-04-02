/**
 * IP 限频工具（内存 Map 实现）
 * 同一 IP 每分钟最多 maxRequests 次请求
 */

interface RateLimitEntry {
  count: number
  resetAt: number // 时间戳 ms
}

const ipMap = new Map<string, RateLimitEntry>()

// 每 5 分钟清理过期条目，防止内存泄漏
setInterval(() => {
  const now = Date.now()
  for (const [ip, entry] of ipMap) {
    if (now > entry.resetAt) {
      ipMap.delete(ip)
    }
  }
}, 5 * 60 * 1000)

/**
 * 检查 IP 是否超出限频
 * @param ip 客户端 IP
 * @param maxRequests 每分钟最大请求数（默认 3）
 * @returns true = 允许通过, false = 超出限制
 */
export function checkRateLimit(ip: string, maxRequests = 3): boolean {
  const now = Date.now()
  const entry = ipMap.get(ip)

  if (!entry || now > entry.resetAt) {
    // 无记录或已过期，重新计数
    ipMap.set(ip, { count: 1, resetAt: now + 60 * 1000 })
    return true
  }

  if (entry.count < maxRequests) {
    entry.count++
    return true
  }

  // 超出限制
  return false
}

/**
 * 从请求事件中提取客户端 IP
 */
export function getClientIp(event: any): string {
  // 优先读 x-forwarded-for（Nginx 反代场景）
  const forwarded = getRequestHeader(event, 'x-forwarded-for')
  if (forwarded) {
    return forwarded.split(',')[0].trim()
  }

  // 其次读 x-real-ip
  const realIp = getRequestHeader(event, 'x-real-ip')
  if (realIp) {
    return realIp.trim()
  }

  // 最后用 socket
  return event.node?.req?.socket?.remoteAddress || 'unknown'
}
