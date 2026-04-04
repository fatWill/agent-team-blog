/**
 * 图片 CDN URL 转换工具
 * 将 /uploads/xxx 相对路径转换为 CDN 完整 URL
 */

const CDN_HOST = 'https://cdn.fatwill.cloud'

/**
 * 将 /uploads/xxx 路径转换为 CDN 完整 URL
 * 已经是完整 URL（http/https 开头）的不做处理
 */
export function toCdnUrl(url: string | null | undefined): string {
  if (!url) return ''
  if (url.startsWith('http')) return url
  if (url.startsWith('/uploads/')) {
    return `${CDN_HOST}${url}`
  }
  return url
}
