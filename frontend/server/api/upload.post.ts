import { readMultipartFormData } from 'h3'
import { extname } from 'node:path'
import { requireAuth } from '~/server/utils/auth'
import { generateCOSKey, uploadToCOS } from '~/server/utils/cos'

/** 允许的图片格式 */
const ALLOWED_EXTENSIONS = new Set(['.jpg', '.jpeg', '.png', '.gif', '.webp'])

/** 允许的 MIME 类型 */
const ALLOWED_MIMES = new Set([
  'image/jpeg',
  'image/png',
  'image/gif',
  'image/webp',
])

/**
 * POST /api/upload
 * 上传图片到腾讯云 COS
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

  // 生成 COS 存储路径并上传
  const key = generateCOSKey(ext)
  const url = await uploadToCOS(key, fileField.data, fileField.type)

  return { url }
})
