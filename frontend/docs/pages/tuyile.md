# 页面：图一乐（/toys/tuyile）

- **文件**: `pages/toys/tuyile.vue`（薄层，`layout: false`）
- **Feature**: `features/tuyile`
- **视图切换**: 用 query 驱动，便于分享和前进后退
  - `/toys/tuyile`：首页入口
  - `/toys/tuyile?view=almanac[&date=YYYY-MM-DD]`：老黄历（不带 date 时为北京时间今天；切日用 `router.replace`）
  - `/toys/tuyile?view=chat`：运势聊天（`<ClientOnly>`，刷新即清空对话）
- **返回**: 左上角链接到 `/toys`
- **SSR**: 黄历通过 `useAsyncData` 在服务端预取，聊天只在客户端渲染
- **SEO**: `useSeoMeta` 设置 title、description、og，ogUrl 通过 `useSiteUrl().pageUrl()` 生成
- **响应式**: 移动端单列；宜/忌、吉神/凶神在 `sm` 及以上为双列；时辰网格为 3/4/6 列
