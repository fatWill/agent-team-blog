import { readMultipartFormData } from 'h3'
import { writeFileSync, mkdirSync, existsSync } from 'node:fs'
import { join } from 'node:path'
import { requireAuth } from '~/server/utils/auth'

/**
 * POST /api/upload/chunk
 * 上传单个分片，保存到临时目录 tmp/{uploadId}/{chunkIndex}
 * 需要鉴权
 */
export default defineEventHandler(async (event) => {
  await requireAuth(event)

  const formData = await readMultipartFormData(event)

  if (!formData || formData.length === 0) {
    throw createError({
      statusCode: 400,
      statusMessage: '请求体为空',
    })
  }

  // 提取表单字段
  const fields: Record<string, string> = {}
  let fileData: Buffer | null = null

  for (const field of formData) {
    if (field.name === 'file') {
      fileData = field.data
    }
    else if (field.name && field.data) {
      fields[field.name] = field.data.toString('utf-8')
    }
  }

  const { uploadId, chunkIndex, totalChunks, filename } = fields

  // 参数校验
  if (!fileData) {
    throw createError({
      statusCode: 400,
      statusMessage: '缺少 file 字段',
    })
  }

  if (!uploadId || chunkIndex === undefined || !totalChunks || !filename) {
    throw createError({
      statusCode: 400,
      statusMessage: '缺少必要参数：uploadId, chunkIndex, totalChunks, filename',
    })
  }

  const idx = Number(chunkIndex)
  const total = Number(totalChunks)

  if (Number.isNaN(idx) || Number.isNaN(total) || idx < 0 || total <= 0 || idx >= total) {
    throw createError({
      statusCode: 400,
      statusMessage: 'chunkIndex 或 totalChunks 参数无效',
    })
  }

  // uploadId 安全校验：仅允许字母、数字、短横线、下划线
  if (!/^[\w-]+$/.test(uploadId)) {
    throw createError({
      statusCode: 400,
      statusMessage: 'uploadId 包含非法字符',
    })
  }

  // 确定临时目录
  // 生产环境：cwd = .output/server，临时目录在 .output/tmp/
  // 开发环境：cwd = 项目根目录，临时目录在 tmp/
  const isProduction = process.env.NODE_ENV === 'production'
  const tmpDir = isProduction
    ? join('/root/blog-uploads/tmp', uploadId)
    : join(process.cwd(), 'tmp', uploadId)

  // 确保目录存在
  if (!existsSync(tmpDir)) {
    mkdirSync(tmpDir, { recursive: true, mode: 0o755 })
  }

  // 写入分片文件
  const chunkPath = join(tmpDir, String(idx))
  writeFileSync(chunkPath, fileData, { mode: 0o644 })

  return {
    received: idx,
    total,
  }
})
