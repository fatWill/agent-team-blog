/**
 * 图片 CDN URL 转换工具
 * 支持新 COS URL（upload/xxx）和旧本地路径（/uploads/xxx）
 * 集成腾讯云数据万象 WebP 自动压缩
 *
 * 域名不再硬编码：assets host 由 `plugins/image-url-host.ts` 在应用启动时
 * 从 runtimeConfig.public.assetsUrl 注入（见 setImageUrlAssetsHost）。
 * 纯计算逻辑通过 createImageUrlHelper 工厂暴露，便于测试和独立复用。
 */

/** 兜底 assets host，与 nuxt.config.ts runtimeConfig.public.assetsUrl 默认值保持一致 */
const DEFAULT_ASSETS_HOST = 'https://assets.fatwill.cloud'

/** 不应追加数据万象参数的文件扩展名 */
const SKIP_EXTENSIONS = ['.svg', '.gif']

export interface ImageUrlHelper {
  toCdnUrl: (url: string | null | undefined) => string
  toThumbUrl: (url: string | null | undefined, width?: number) => string
  toWebpUrl: (url: string | null | undefined) => string
}

/**
 * 创建一组绑定了指定 assets host 的图片 URL 转换函数（纯函数工厂，无副作用）
 * @param assetsHost 图片 CDN 域名，例如 `https://assets.example.com`（末尾斜杠会被移除）
 */
export function createImageUrlHelper(assetsHost: string): ImageUrlHelper {
  const host = assetsHost.replace(/\/+$/, '')

  /**
   * 判断 URL 是否可以追加数据万象参数
   * - 必须属于当前 assets CDN 域名
   * - 不能已带 imageMogr2 参数
   * - 不能是 SVG / GIF 格式
   */
  function canApplyImageMogr(cdnUrl: string): boolean {
    if (!cdnUrl.startsWith(host)) return false
    if (cdnUrl.includes('imageMogr2')) return false
    const lower = cdnUrl.toLowerCase()
    return !SKIP_EXTENSIONS.some(ext => lower.includes(ext))
  }

  /**
   * 将图片路径转换为完整的 CDN URL
   * - 已经是完整 URL（http/https 开头）→ 直接返回
   * - 旧格式 `/uploads/xxx` → `<assetsHost>/uploads/xxx`（兼容旧数据）
   * - 新格式 `upload/xxx` → `<assetsHost>/upload/xxx`
   */
  function toCdnUrl(url: string | null | undefined): string {
    if (!url) return ''
    // 已经是完整 URL，直接返回
    if (url.startsWith('http://') || url.startsWith('https://')) return url
    // 旧格式：/uploads/xxx（带前导斜杠）
    if (url.startsWith('/uploads/')) {
      return `${host}${url}`
    }
    // 新格式：upload/xxx（无前导斜杠）
    if (url.startsWith('upload/')) {
      return `${host}/${url}`
    }
    return url
  }

  /**
   * 将图片路径转换为数据万象缩略图 URL（列表页封面图）
   * 参数：800px 宽度限制 + WebP + 80% 质量 + 去 EXIF
   * @param url 原始图片路径
   * @param width 缩略图最大宽度，默认 800
   */
  function toThumbUrl(url: string | null | undefined, width = 800): string {
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
  function toWebpUrl(url: string | null | undefined): string {
    const cdnUrl = toCdnUrl(url)
    if (!cdnUrl) return ''
    if (!canApplyImageMogr(cdnUrl)) return cdnUrl
    return `${cdnUrl}?imageMogr2/format/webp/rquality/85/strip`
  }

  return { toCdnUrl, toThumbUrl, toWebpUrl }
}

/** 当前生效的 helper 实例（默认用兜底 host，启动插件会用 runtimeConfig 值替换） */
let helper: ImageUrlHelper = createImageUrlHelper(DEFAULT_ASSETS_HOST)

/**
 * 注入运行时 assets host（由 plugins/image-url-host.ts 在 SSR 与客户端启动时调用）
 * @param assetsHost runtimeConfig.public.assetsUrl
 */
export function setImageUrlAssetsHost(assetsHost: string): void {
  if (!assetsHost) return
  helper = createImageUrlHelper(assetsHost)
}

// ====== 对外 API（签名与改造前完全一致，调用方无需修改） ======

/** 将图片路径转换为完整的 CDN URL */
export function toCdnUrl(url: string | null | undefined): string {
  return helper.toCdnUrl(url)
}

/** 将图片路径转换为数据万象缩略图 URL（列表页封面图） */
export function toThumbUrl(url: string | null | undefined, width = 800): string {
  return helper.toThumbUrl(url, width)
}

/** 将图片路径转换为 WebP 大图 URL（文章详情页内图片） */
export function toWebpUrl(url: string | null | undefined): string {
  return helper.toWebpUrl(url)
}
