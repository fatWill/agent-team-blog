# 微信 JS-SDK 接口文档

## [GET] /api/wechat/jssdk-config - 获取 JS-SDK 签名配置

### 描述

获取微信 JS-SDK 的签名配置信息，供前端调用 `wx.config()` 使用，实现自定义分享卡片等功能。

### 鉴权

❌ 不需要（公开接口）

### 请求参数（Query）

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| url | string | ✅ | 当前页面的完整 URL（含协议、域名、路径、query，不含 hash） |

### 请求示例

```
GET /api/wechat/jssdk-config?url=https://fatwill.cloud/articles/d76f4f63-9c98-4135-af90-7dfcb05c2cd8
```

### 成功响应 (200)

```json
{
  "appId": "wx13c03e11e4bf5469",
  "timestamp": 1747126644,
  "nonceStr": "aBcDeFgH12345678",
  "signature": "e10adc3949ba59abbe56e057f20f883e"
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| appId | string | 微信公众号 AppID |
| timestamp | int64 | 签名时间戳（秒） |
| nonceStr | string | 随机字符串（16 位） |
| signature | string | SHA1 签名 |

### 错误响应

| statusCode | statusMessage | 场景 |
|------------|---------------|------|
| 400 | 缺少 url 参数 | 未传入 url 参数 |
| 500 | 微信公众号未配置 | 环境变量 WECHAT_APP_ID / WECHAT_APP_SECRET 未设置 |
| 500 | 获取微信签名失败 | 调用微信接口失败（网络异常、IP 白名单未配置等） |

### 前端使用示例

```javascript
// 1. 调用后端接口获取签名
const res = await fetch(`/api/wechat/jssdk-config?url=${encodeURIComponent(location.href.split('#')[0])}`);
const config = await res.json();

// 2. 调用 wx.config 注入配置
wx.config({
  debug: false,
  appId: config.appId,
  timestamp: config.timestamp,
  nonceStr: config.nonceStr,
  signature: config.signature,
  jsApiList: ['updateAppMessageShareData', 'updateTimelineShareData']
});

// 3. 自定义分享内容
wx.ready(() => {
  wx.updateAppMessageShareData({
    title: '文章标题',
    desc: '文章摘要',
    link: window.location.href,
    imgUrl: 'https://assets.fatwill.cloud/cover.jpg',
    success: () => {}
  });
});
```

### 实现细节

- **access_token 缓存**：Redis Key `wechat:access_token`，TTL 7000s（微信有效期 7200s，留 200s 安全余量）
- **jsapi_ticket 缓存**：Redis Key `wechat:jsapi_ticket`，TTL 7000s
- **签名算法**：`SHA1(jsapi_ticket=xxx&noncestr=xxx&timestamp=xxx&url=xxx)`，字段按字典序排列
- **Token 过期处理**：当 jsapi_ticket 接口返回 40001/42001 错误码时，自动清除 access_token 缓存
