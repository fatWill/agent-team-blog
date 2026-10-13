import { rememberFirstEntryUrl } from '~/features/tuyile/composables/useTuyileWechat'

/**
 * 记录站点首次进入的 URL（去掉 #）。
 * iOS 微信对 SPA pushState 以首次进入地址为准做 JS-SDK 签名，图一乐分享失败时用它重签。
 */
export default defineNuxtPlugin(() => {
  rememberFirstEntryUrl()
})
