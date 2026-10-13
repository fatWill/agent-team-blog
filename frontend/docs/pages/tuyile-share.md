# 页面：图一乐对话分享（/toys/tuyile/share/:hash）

- **文件**: `pages/toys/tuyile/share/[hash].vue`（`layout: false`，外壳 `TuyileShell`）
- **数据**: SSR `useAsyncData` → `GET /api/toys/tuyile/share/:hash`；hash 非 12 位 base62 时前端/代理直接 404
- **只读**: 只渲染消息（`ChatMessage`，markdown-it `html:false`；user 纯文本插值），无输入框、无「继续聊」入口；底部「我也去算一卦」跳 `/toys/tuyile?view=chat`（全新对话）
- **404**: 友好页「这一卦找不到啦」，SSR 返回真实 404 状态码
- **SEO/OG**: title 取首条 user 摘要（24 字）；description / og:description 取首条 assistant 摘要（100 字，Markdown 转纯文本，尖括号转全角）；og:image 用站点 512 图标；`robots: noindex, follow`（用户生成内容不收录）
- **微信**: 打开即设置分享卡片（link 为本页），便于二次转发；hash 变化后在新 URL 重签，失败时用首次进入 URL 重签（iOS）
- **缓存**: 后端 `Cache-Control: public, max-age=3600` 仅作用于 `/api/.../share/:hash` JSON；页面 HTML 由 Nuxt SSR 输出，未设置页面级缓存
