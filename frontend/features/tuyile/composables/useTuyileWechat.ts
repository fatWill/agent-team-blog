import { useWechatShare } from '~/composables/useWechatShare'

/**
 * 图一乐微信分享卡片
 * - 每次调用都用「当前地址（去掉 # 及之后，保留 query）」重新签名，适配 SPA 内路由切换
 * - iOS 微信可能以首次进入的地址为准签名：当前地址签名失败时，用首次进入 URL 再签一次
 */

let firstEntryUrl = ''

/** 记录首次进入的 URL（在客户端插件/页面最早时机调用，只记一次） */
export function rememberFirstEntryUrl() {
  if (import.meta.server || firstEntryUrl) return
  firstEntryUrl = window.location.href.split('#')[0]
}

export interface TuyileCard {
  title: string
  desc: string
  link: string
  imgUrl: string
}

export async function setTuyileWechatCard(card: TuyileCard): Promise<boolean> {
  const { isWechat, initWechatShare } = useWechatShare()
  if (!isWechat()) return false
  rememberFirstEntryUrl()
  const current = window.location.href.split('#')[0]
  const ok = await initWechatShare(card, { signUrl: current })
  if (ok || !firstEntryUrl || firstEntryUrl === current) return ok
  return initWechatShare(card, { signUrl: firstEntryUrl })
}
