import { getDb } from '~/server/utils/db'

/**
 * IP 限频：同一 IP 每分钟最多 3 次请求
 * 内存 Map 实现（Nuxt 单进程，不需要 Redis）
 */
const ipLimitMap = new Map<string, { count: number; resetAt: number }>()

function checkIpLimit(ip: string): boolean {
  const now = Date.now()
  const record = ipLimitMap.get(ip)

  if (!record || now >= record.resetAt) {
    ipLimitMap.set(ip, { count: 1, resetAt: now + 60 * 1000 })
    return true
  }

  if (record.count >= 3) {
    return false
  }

  record.count++
  return true
}

// 定期清理过期记录，防止内存泄漏（每 5 分钟清理一次）
setInterval(() => {
  const now = Date.now()
  for (const [ip, record] of ipLimitMap) {
    if (now >= record.resetAt) {
      ipLimitMap.delete(ip)
    }
  }
}, 5 * 60 * 1000)

/**
 * 获取客户端 IP
 */
function getClientIp(event: any): string {
  const xff = getRequestHeader(event, 'x-forwarded-for')
  if (xff) {
    return xff.split(',')[0].trim()
  }
  return event.node?.req?.socket?.remoteAddress || 'unknown'
}

/**
 * POST /api/messages
 * 新增留言（公开接口）
 * Body: { deviceId: string, nickname?: string, content: string }
 */
export default defineEventHandler(async (event) => {
  // IP 限频检查
  const clientIp = getClientIp(event)
  if (!checkIpLimit(clientIp)) {
    throw createError({
      statusCode: 429,
      statusMessage: '操作太频繁，请稍后再试',
    })
  }

  const body = await readBody<{ deviceId: string; nickname?: string; content: string }>(event)

  const deviceId = body?.deviceId?.trim()
  if (!deviceId || deviceId.length < 8 || deviceId.length > 64) {
    throw createError({ statusCode: 400, statusMessage: '无效的设备 ID' })
  }

  const content = body?.content?.trim()
  if (!content) {
    throw createError({ statusCode: 400, statusMessage: '留言内容不能为空' })
  }

  const nickname = body?.nickname?.trim() || null

  const db = getDb()

  // 检查该设备是否已有留言
  const existRow = db.prepare(
    'SELECT id, last_modified_date FROM messages WHERE device_id = ? LIMIT 1',
  ).get(deviceId) as { id: number; last_modified_date: string | null } | undefined

  // 计算今天的日期字符串（UTC+8）
  const now = new Date()
  const utc8 = new Date(now.getTime() + 8 * 60 * 60 * 1000)
  const todayStr = utc8.toISOString().slice(0, 10)

  if (existRow) {
    // 已有留言 → 视为修改
    const lastModDate = existRow.last_modified_date ? String(existRow.last_modified_date).slice(0, 10) : null

    if (lastModDate === todayStr) {
      throw createError({
        statusCode: 403,
        statusMessage: '今天已经修改过了，明天再来吧',
      })
    }

    // 更新留言
    db.prepare(
      'UPDATE messages SET nickname = ?, content = ?, last_modified_date = ?, ip = ? WHERE id = ?',
    ).run(nickname, content, todayStr, clientIp, existRow.id)

    // 查询更新后的记录
    const row = db.prepare(
      'SELECT id, nickname, content, last_modified_date, created_at, updated_at FROM messages WHERE id = ?',
    ).get(existRow.id) as any

    return {
      id: row.id,
      nickname: row.nickname || '匿名',
      content: row.content,
      isOwn: true,
      canEdit: false,
      createdAt: row.created_at,
      updatedAt: row.updated_at,
    }
  }
  else {
    // 新增留言（last_modified_date 为 NULL，不算修改过）
    const result = db.prepare(
      'INSERT INTO messages (device_id, nickname, content, ip) VALUES (?, ?, ?, ?)',
    ).run(deviceId, nickname, content, clientIp)

    // 查询新插入的记录
    const row = db.prepare(
      'SELECT id, nickname, content, last_modified_date, created_at, updated_at FROM messages WHERE id = ?',
    ).get(Number(result.lastInsertRowid)) as any

    return {
      id: row.id,
      nickname: row.nickname || '匿名',
      content: row.content,
      isOwn: true,
      canEdit: true,
      createdAt: row.created_at,
      updatedAt: row.updated_at,
    }
  }
})
