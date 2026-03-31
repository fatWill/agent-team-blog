import { readMultipartFormData } from 'h3'
import { writeFileSync, mkdirSync, existsSync } from 'node:fs'
import { join, extname } from 'node:path'
import { requireAuth } from '~/server/utils/auth'

/** 允许的图片格式 */
const ALLOWED_EXTENSIONS = new Set(['.jpg', '.jpeg', '.png', '.gif', '.webp'])

/** 允许的 MIME 类型 */
const ALLOWED_MIMES = new Set([
  'image/jpeg',
  'image/png',
  'image/gif',
  'image/webp',
])

/** 最大文件大小：5MB */
const MAX_FILE_SIZE = 5 * 1024 * 1024

/**
 * 生成随机字符串
 */
function randomString(len: number): string {
  const chars = 'abcdefghijklmnopqrstuvwxyz0123456789'
  let result = ''
  for (let i = 0; i < len; i++) {
    result += chars.charAt(Math.floor(Math.random() * chars.length))
  }
  return result
}

/**
 * POST /api/upload
 * 上传图片，保存到 public/uploads/ 目录
 * 需要鉴权
 */
export default defineEventHandler(async (event) => {
  // 鉴权校验
  await requireAuth(event)

  // 读取 multipart 表单数据
  const formData = await readMultipartFormData(event)

  if (!formData || formData.length === 0) {
    throw createError({
      statusCode: 400,
      statusMessage: '请选择要上传的文件',
    })
  }

  // 查找 file 字段
  const fileField = formData.find(f => f.name === 'file')

  if (!fileField || !fileField.data || !fileField.filename) {
    throw createError({
      statusCode: 400,
      statusMessage: '缺少 file 字段或文件为空',
    })
  }

  // 校验文件大小
  if (fileField.data.length > MAX_FILE_SIZE) {
    throw createError({
      statusCode: 400,
      statusMessage: '文件大小不能超过 5MB',
    })
  }

  // 校验 MIME 类型
  if (fileField.type && !ALLOWED_MIMES.has(fileField.type)) {
    throw createError({
      statusCode: 400,
      statusMessage: '不支持的文件格式，仅支持 jpg/jpeg/png/gif/webp',
    })
  }

  // 校验文件扩展名
  const ext = extname(fileField.filename).toLowerCase()
  if (!ALLOWED_EXTENSIONS.has(ext)) {
    throw createError({
      statusCode: 400,
      statusMessage: '不支持的文件格式，仅支持 jpg/jpeg/png/gif/webp',
    })
  }

  // 生成唯一文件名：时间戳-随机串.扩展名
  const filename = `${Date.now()}-${randomString(8)}${ext}`

  // 确定存储目录
  // 生产环境：cwd = .output/server，静态资源在 .output/public/uploads/
  // 开发环境：cwd = 项目根目录，静态资源在 public/uploads/
  const isProduction = process.env.NODE_ENV === 'production'
  const uploadDir = isProduction
    ? join(process.cwd(), '..', 'public', 'uploads')
    : join(process.cwd(), 'public', 'uploads')

  // 确保目录存在（mode 0o755 确保 Nginx 可访问）
  if (!existsSync(uploadDir)) {
    mkdirSync(uploadDir, { recursive: true, mode: 0o755 })
  }

  // 写入文件（mode 0o644 确保 Nginx worker 有读权限）
  const filePath = join(uploadDir, filename)
  writeFileSync(filePath, fileField.data, { mode: 0o644 })

  // 返回可访问的 URL
  return {
    url: `/uploads/${filename}`,
  }
})
