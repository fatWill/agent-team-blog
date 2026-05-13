/**
 * 微信 JS-SDK 分享 composable
 * 仅在微信内置浏览器中生效，自动调用后端获取签名并配置分享卡片
 */

interface WechatShareConfig {
  /** 分享标题 */
  title: string
  /** 分享描述 */
  desc: string
  /** 分享链接（当前页面 URL，不含 # 及后面的内容） */
  link: string
  /** 分享缩略图 URL（建议 300x300 以上正方形） */
  imgUrl: string
}

interface JsSdkConfigResponse {
  appId: string
  timestamp: number
  nonceStr: string
  signature: string
}

/** 判断是否在微信内置浏览器中 */
function isWechat(): boolean {
  if (import.meta.server) return false
  return /MicroMessenger/i.test(navigator.userAgent)
}

/**
 * 初始化微信 JS-SDK 分享配置
 * - 仅在微信内置浏览器中执行
 * - 自动调用后端获取 JS-SDK 签名
 * - 配置「发送给朋友」和「分享到朋友圈」的自定义卡片
 */
async function initWechatShare(config: WechatShareConfig): Promise<void> {
  // 1. 非微信浏览器直接返回
  if (!isWechat()) return

  try {
    // 2. 获取当前页面 URL（不含 # 及后面的内容），用于签名
    const signUrl = window.location.href.split('#')[0]

    // 3. 调用后端接口获取 JS-SDK 签名配置
    const sdkConfig = await $fetch<JsSdkConfigResponse>('/api/wechat/jssdk-config', {
      params: { url: signUrl },
    })

    // 4. 动态 import weixin-js-sdk（避免 SSR 报错）
    const wx = (await import('weixin-js-sdk')).default

    // 5. 配置 wx.config
    wx.config({
      debug: false,
      appId: sdkConfig.appId,
      timestamp: sdkConfig.timestamp,
      nonceStr: sdkConfig.nonceStr,
      signature: sdkConfig.signature,
      jsApiList: ['updateAppMessageShareData', 'updateTimelineShareData'],
    })

    // 6. wx.ready 中设置分享数据
    wx.ready(() => {
      // 分享给朋友
      wx.updateAppMessageShareData({
        title: config.title,
        desc: config.desc,
        link: config.link,
        imgUrl: config.imgUrl,
        success() {
          // 设置成功，用户点击分享时会使用这些数据
        },
      })

      // 分享到朋友圈
      wx.updateTimelineShareData({
        title: config.title,
        link: config.link,
        imgUrl: config.imgUrl,
        success() {
          // 设置成功
        },
      })
    })

    // 7. 错误处理
    wx.error((res: { errMsg: string }) => {
      console.error('[WeChat JS-SDK] wx.config error:', res)
    })
  } catch (err) {
    console.error('[WeChat JS-SDK] initWechatShare failed:', err)
  }
}

export function useWechatShare() {
  return {
    isWechat,
    initWechatShare,
  }
}
