import { getProfile } from '~/server/utils/profile'

/**
 * GET /api/profile
 * 获取博主个人资料（公开接口，无需鉴权）
 */
export default defineEventHandler(async () => {
  return await getProfile()
})
