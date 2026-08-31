/**
 * 站点域名 / URL 拼接 composable
 *
 * 所有绝对 URL（canonical、OG、分享链接、CDN 素材）统一从此处获取，
 * 避免域名硬编码。默认值定义在 nuxt.config.ts 的 runtimeConfig.public，
 * 可通过环境变量覆盖：
 *   - NUXT_PUBLIC_SITE_URL    站点主域名
 *   - NUXT_PUBLIC_ASSETS_URL  图片/素材 CDN
 *   - NUXT_PUBLIC_CDN_URL     构建产物 CDN
 */
export const useSiteUrl = () => {
  const config = useRuntimeConfig()

  const siteUrl = config.public.siteUrl
  const assetsUrl = config.public.assetsUrl
  const cdnUrl = config.public.cdnUrl

  return {
    siteUrl,
    assetsUrl,
    cdnUrl,

    /** 拼接文章详情页绝对 URL */
    articleUrl: (id: string) => `${siteUrl}/articles/${id}`,

    /** 拼接任意页面的绝对 URL（自动补前导斜杠） */
    pageUrl: (path: string) => `${siteUrl}${path.startsWith('/') ? path : `/${path}`}`,

    /** OG 默认兜底图 */
    defaultOgImage: () => `${siteUrl}/og-default.png`,

    /** 拼接 assets CDN 素材绝对 URL（自动补前导斜杠） */
    assetUrl: (path: string) => `${assetsUrl}${path.startsWith('/') ? path : `/${path}`}`,
  }
}
