import { readBody } from 'h3'
import { readFileSync, existsSync, rmSync } from 'node:fs'
import { join, extname } from 'node:path'
import { requireAuth } from '~/server/utils/auth'
import { generateCOSKey, uploadToCOS } from '~/server/utils/cos'

/** 允许的文件扩展名 */
const ALLOWED_EXTENSIONS = new Set(['.jpg', '.jpeg', '.png', '.gif', '.webp'])

/** MIME 类型映射 */
const MIME_MAP: Record<string, string> = {
  '.jpg': 'image/jpeg',
  '.jpeg': 'image/jpeg',
  '.png': 'image/png',
  '.gif': 'image/gif',
  '.webp': 'image/webp',
}

/**
 * POST /api/upload/merge
 * 合并分片并上传到腾讯云 COS，清理临时目录
 * 需要鉴权
 */
export default defineEventHandler(async (event) => {
  await requireAuth(event)

  const body = await readBody<{
    uploadId: string
    totalChunks: number
    filename: string
  }>(event)

  if (!body?.uploadId || !body?.totalChunks || !body?.filename) {
    throw createError({
      statusCode: 400,
      statusMessage: '缺少必要参数：uploadId, totalChunks, filename',
    })
  }

  const { uploadId, totalChunks, filename } = body

  // uploadId 安全校验
  if (!/^[\w-]+$/.test(uploadId)) {
    throw createError({
      statusCode: 400,
      statusMessage: 'uploadId 包含非法字符',
    })
  }

  // 校验文件扩展名
  const ext = extname(filename).toLowerCase()
  if (!ALLOWED_EXTENSIONS.has(ext)) {
    throw createError({
      statusCode: 400,
      statusMessage: '不支持的文件格式，仅支持 jpg/jpeg/png/gif/webp',
    })
  }

  const isProduction = process.env.NODE_ENV === 'production'
  const tmpDir = isProduction
    ? join('/root/blog-uploads/tmp', uploadId)
    : join(process.cwd(), 'tmp', uploadId)

  // 校验临时目录是否存在
  if (!existsSync(tmpDir)) {
    throw createError({
      statusCode: 400,
      statusMessage: '上传会话不存在或已过期',
    })
  }

  // 按序读取所有分片并合并
  const chunks: Buffer[] = []
  for (let i = 0; i < totalChunks; i++) {
    const chunkPath = join(tmpDir, String(i))
    if (!existsSync(chunkPath)) {
      throw createError({
        statusCode: 400,
        statusMessage: `分片 ${i} 缺失，请重新上传`,
      })
    }
    chunks.push(readFileSync(chunkPath))
  }

  const mergedBuffer = Buffer.concat(chunks)

  // 生成 COS 存储路径并上传
  const key = generateCOSKey(ext)
  const contentType = MIME_MAP[ext] || 'application/octet-stream'
  const url = await uploadToCOS(key, mergedBuffer, contentType)

  // 清理临时目录
  rmSync(tmpDir, { recursive: true, force: true })

  return { url }
})
