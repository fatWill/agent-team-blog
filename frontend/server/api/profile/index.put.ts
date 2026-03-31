import { readBody } from 'h3'
import { requireAuth } from '~/server/utils/auth'
import { updateProfile } from '~/server/utils/profile'

/**
 * PUT /api/profile
 * 更新博主个人资料（需鉴权）
 */
export default defineEventHandler(async (event) => {
  // 鉴权校验
  await requireAuth(event)

  const body = await readBody<{
    avatar?: string
    bio?: string
  }>(event)

  // 至少需要提供一个更新字段
  if (body?.avatar === undefined && body?.bio === undefined) {
    throw createError({
      statusCode: 400,
      statusMessage: '至少需要提供 avatar 或 bio 中的一个字段',
    })
  }

  // 构建更新数据
  const updates: { avatar?: string; bio?: string } = {}

  if (body.avatar !== undefined) {
    updates.avatar = typeof body.avatar === 'string' ? body.avatar.trim() : ''
  }

  if (body.bio !== undefined) {
    updates.bio = typeof body.bio === 'string' ? body.bio.trim() : ''
  }

  return await updateProfile(updates)
})
