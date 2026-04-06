/**
 * 图片 CDN URL 转换工具
 * 支持新 COS URL（upload/xxx）和旧本地路径（/uploads/xxx）
 */

const ASSETS_HOST = 'https://assets.fatwill.cloud'

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
 * 将图片路径转换为数据万象缩略图 URL
 * 仅对 assets.fatwill.cloud 域名的图片追加缩略图参数
 * 旧格式本地路径（/uploads/xxx）不做处理，直接走原图
 * @param url 原始图片路径
 * @param size 缩略图最大边长，默认 400
 */
export function toThumbUrl(url: string | null | undefined, size = 400): string {
  const cdnUrl = toCdnUrl(url)
  if (!cdnUrl) return ''
  // 只对 COS/CDN 域名的图片追加缩略图参数，旧本地路径不处理
  if (!cdnUrl.includes('assets.fatwill.cloud')) return cdnUrl
  return `${cdnUrl}?imageMogr2/thumbnail/${size}x${size}>/format/webp/rquality/75/strip`
}
