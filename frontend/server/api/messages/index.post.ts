import { checkRateLimit, getClientIp } from '~/server/utils/rateLimit'
import {
  getMessageByDeviceId,
  createMessage,
  updateMessage,
  getTodayDateCST,
} from '~/server/utils/messages'

/**
 * POST /api/messages
 * 新增留言（公开接口）
 * Body: { deviceId: string, nickname?: string, content: string }
 * - IP 限频：同一 IP 每分钟最多 3 次
 * - 如果 device_id 已有留言，视为修改
 */
export default defineEventHandler(async (event) => {
  // IP 限频检查
  const ip = getClientIp(event)
  if (!checkRateLimit(ip, 3)) {
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

  const nickname = body?.nickname?.trim() || undefined

  // 检查该设备是否已有留言
  const existing = await getMessageByDeviceId(deviceId)

  if (existing) {
    // 已有留言 → 视为修改，检查每日限制
    const today = getTodayDateCST()
    if (existing.last_modified_date === today) {
      throw createError({
        statusCode: 403,
        statusMessage: '今天已经修改过了，明天再来吧',
      })
    }

    const message = await updateMessage(existing.id, {
      nickname,
      content,
      deviceId,
    })
    return message
  }

  // 新增留言
  const message = await createMessage({
    deviceId,
    nickname,
    content,
    ip,
  })

  return message
})
