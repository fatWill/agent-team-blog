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
