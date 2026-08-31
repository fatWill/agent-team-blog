/**
 * 注入图片 CDN 域名到 imageUrl 工具
 *
 * `shared/utils/imageUrl.ts` 是纯工具模块，无法直接使用 composable，
 * 因此在应用启动最早期（文件名 `0.` 前缀保证优先级）把 runtimeConfig 中的
 * assetsUrl 注入进去，之后所有 toCdnUrl / toThumbUrl / toWebpUrl 调用都会
 * 使用运行时配置的域名。
 *
 * 注意：SSR 与客户端都需要执行，故不加 .client / .server 后缀。
 */
import { setImageUrlAssetsHost } from '~/shared/utils/imageUrl'

export default defineNuxtPlugin(() => {
  const config = useRuntimeConfig()
  setImageUrlAssetsHost(config.public.assetsUrl)
})
