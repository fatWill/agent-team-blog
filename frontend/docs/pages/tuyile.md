# 页面：图一乐（/toys/tuyile）

- **文件**: `pages/toys/tuyile/index.vue`（薄层，`layout: false`，外壳 `TuyileShell`）。二期由 `pages/toys/tuyile.vue` 迁为目录路由，便于挂子路由 `share/[hash]`；路由名仍为 `toys-tuyile`，`/toys/tuyile?view=...` 不受影响
- **Feature**: `features/tuyile`
- **视图切换**（query 驱动）
  - `/toys/tuyile`：首页（顶部运势 banner + 3 入口）
  - `/toys/tuyile?view=almanac[&date=YYYY-MM-DD]`：老黄历（顶部 banner；日期条：前后翻日 / 月历 / 今天 / 快速选日）
  - `/toys/tuyile?view=chat`：运势聊天（`<ClientOnly>`，刷新即清空对话）
- **快速选日**：antd `DatePicker`（日/月/年/十年面板），`defineAsyncComponent` 首次点击才加载，包在 `<ClientOnly>`；中文 locale（`ant-design-vue/es/locale/zh_CN` + `dayjs/locale/zh-cn`）；`disabledDate` 限制 1900-01-01 ~ 2100-12-31；主题色跟随暗色模式
- **返回**: 左上角 `/toys`
- **SSR**: 黄历 `useAsyncData` 服务端预取；聊天仅客户端
- **SEO**: `useSeoMeta`；ogUrl 走 `useSiteUrl().pageUrl()`
- **微信**: 进入及 query 变化后，在当前 URL（去 #，保留 query）上重签，默认卡片指向 `/toys/tuyile`
