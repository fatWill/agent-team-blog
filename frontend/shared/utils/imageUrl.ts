/**
 * 图片 CDN URL 转换工具
 * 支持新 COS URL（upload/xxx）和旧本地路径（/uploads/xxx）
 * 集成腾讯云数据万象 WebP 自动压缩
 */

const ASSETS_HOST = 'https://assets.fatwill.cloud'

/** 不应追加数据万象参数的文件扩展名 */
const SKIP_EXTENSIONS = ['.svg', '.gif']

/**
 * 判断 URL 是否可以追加数据万象参数
 * - 必须是 assets.fatwill.cloud 域名
 * - 不能已带 imageMogr2 参数
 * - 不能是 SVG / GIF 格式
 */
function canApplyImageMogr(cdnUrl: string): boolean {
  if (!cdnUrl.includes('assets.fatwill.cloud')) return false
  if (cdnUrl.includes('imageMogr2')) return false
  const lower = cdnUrl.toLowerCase()
  return !SKIP_EXTENSIONS.some(ext => lower.includes(ext))
}

/**
 * 将图片路径转换为完整的 CDN URL
 * - 已经是完整 URL（http/https 开头）→ 直接返回
 * - 旧格式 /uploads/xxx → https://assets.fatwill.cloud/uploads/xxx（兼容旧数据）
 * - 新格式 upload/xxx → https://assets.fatwill.cloud/upload/xxx
 */
export function toCdnUrl(url: string | null | undefined): string {
  if (!url) return ''
  // 已经是完整 URL，直接返回
  if (url.startsWith('http://') || url.startsWith('https://')) return url
  // 旧格式：/uploads/xxx（带前导斜杠）
  if (url.startsWith('/uploads/')) {
    return `${ASSETS_HOST}${url}`
  }
  // 新格式：upload/xxx（无前导斜杠）
  if (url.startsWith('upload/')) {
    return `${ASSETS_HOST}/${url}`
  }
  return url
}

/**
 * 将图片路径转换为数据万象缩略图 URL（列表页封面图）
 * 参数：800px 宽度限制 + WebP + 80% 质量 + 去 EXIF
 * @param url 原始图片路径
 * @param width 缩略图最大宽度，默认 800
 */
export function toThumbUrl(url: string | null | undefined, width = 800): string {
  const cdnUrl = toCdnUrl(url)
  if (!cdnUrl) return ''
  if (!canApplyImageMogr(cdnUrl)) return cdnUrl
  return `${cdnUrl}?imageMogr2/thumbnail/${width}x>/format/webp/rquality/80/strip`
}

/**
 * 将图片路径转换为 WebP 大图 URL（文章详情页内图片）
 * 保持原始尺寸，仅转 WebP + 85% 质量 + 去 EXIF
 * @param url 原始图片路径
 */
export function toWebpUrl(url: string | null | undefined): string {
  const cdnUrl = toCdnUrl(url)
  if (!cdnUrl) return ''
  if (!canApplyImageMogr(cdnUrl)) return cdnUrl
  return `${cdnUrl}?imageMogr2/format/webp/rquality/85/strip`
}
