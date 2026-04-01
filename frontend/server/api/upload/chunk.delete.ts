import { readBody } from 'h3'
import { existsSync, rmSync } from 'node:fs'
import { join } from 'node:path'
import { requireAuth } from '~/server/utils/auth'

/**
 * DELETE /api/upload/chunk
 * 取消上传，清理临时分片文件
 * 需要鉴权
 */
export default defineEventHandler(async (event) => {
  await requireAuth(event)

  const body = await readBody<{ uploadId: string }>(event)

  if (!body?.uploadId) {
    throw createError({
      statusCode: 400,
      statusMessage: '缺少必要参数：uploadId',
    })
  }

  const { uploadId } = body

  // uploadId 安全校验：仅允许字母、数字、短横线、下划线
  if (!/^[\w-]+$/.test(uploadId)) {
    throw createError({
      statusCode: 400,
      statusMessage: 'uploadId 包含非法字符',
    })
  }

  // 确定临时目录
  const isProduction = process.env.NODE_ENV === 'production'
  const tmpDir = isProduction
    ? join('/root/blog-uploads/tmp', uploadId)
    : join(process.cwd(), 'tmp', uploadId)

  // 删除临时目录（如果存在）
  if (existsSync(tmpDir)) {
    rmSync(tmpDir, { recursive: true, force: true })
  }

  return { ok: true }
})
