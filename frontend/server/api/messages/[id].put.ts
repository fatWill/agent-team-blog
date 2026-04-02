import { checkRateLimit, getClientIp } from '~/server/utils/rateLimit'
import {
  getMessageById,
  updateMessage,
  getTodayDateCST,
} from '~/server/utils/messages'

/**
 * PUT /api/messages/:id
 * 修改留言（公开接口）
 * Body: { deviceId: string, nickname?: string, content: string }
 * - 验证 device_id 归属
 * - 每天只能修改一次
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

  const id = Number(event.context.params?.id)
  if (!id || isNaN(id)) {
    throw createError({ statusCode: 400, statusMessage: '无效的留言 ID' })
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

  // 查找留言
  const existing = await getMessageById(id)
  if (!existing) {
    throw createError({ statusCode: 404, statusMessage: '留言不存在' })
  }

  // 验证归属
  if (existing.device_id !== deviceId) {
    throw createError({ statusCode: 403, statusMessage: '无权修改此留言' })
  }

  // 检查每日修改限制
  const today = getTodayDateCST()
  if (existing.last_modified_date === today) {
    throw createError({
      statusCode: 403,
      statusMessage: '今天已经修改过了，明天再来吧',
    })
  }

  // 执行更新
  const message = await updateMessage(id, {
    nickname,
    content,
    deviceId,
  })

  return message
})
