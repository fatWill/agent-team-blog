# 图一乐 - 组件清单

## TuyileShell

- **路径**: `features/tuyile/components/TuyileShell.vue`
- **职责**: 页面外壳，包含返回 `/toys`、主题切换，并引入公共样式 `styles/tuyile.css`
- **Slots**: `default`、`nav`（额外悬浮导航）

## FortuneBanner

- **路径**: `features/tuyile/components/FortuneBanner.vue`
- **职责**: 顶部「聊一下你的运势」广告 banner，带流动渐变、扫光和图标轻浮动；`prefers-reduced-motion` 时关闭动效
- **Emits**: `click()`

## TuyileHome

- **路径**: `features/tuyile/components/TuyileHome.vue`
- **职责**: 顶部 banner 加 3 块入口（星座、八字显示「即将上线」）
- **Emits**: `open('almanac')`、`chat()`

## AlmanacView

- **路径**: `features/tuyile/components/AlmanacView.vue`
- **职责**: 老黄历主视图，包括顶部 banner、日期条（前后翻日 / 月历 / 今天 / 快速选日）和全部字段
- **Props**: `date: string`、`today: string`
- **Emits**: `update:date(string)`、`chat()`

## AlmanacDatePicker

- **路径**: `features/tuyile/components/AlmanacDatePicker.vue`（由 AlmanacView 异步加载，放在 `<ClientOnly>` 中）
- **职责**: antd DatePicker 弹层（日 / 月 / 年 / 十年面板），中文 locale，日期限定 1900~2100；输入框隐藏
- **Props**: `modelValue: string`、`open: boolean`、`dark: boolean`
- **Emits**: `update:modelValue(string)`、`update:open(boolean)`

## AlmanacCalendar / AlmanacHours

同一期，见 git 历史。

## FortuneChat

- **路径**: `features/tuyile/components/FortuneChat.vue`
- **职责**: AG-UI 聊天，包括思考状态条、快捷提问 chip、轮换 placeholder、分享按钮（至少 1 条完成的 assistant 回复才可点，生成中禁用）、停止生成和错误提示
- **Emits**: `back()`

## ChatMessage

- **路径**: `features/tuyile/components/ChatMessage.vue`
- **职责**: 单条消息气泡（聊天页与分享页共用），安全渲染 Markdown
- **Props**: `role: 'user' | 'assistant'`、`content: string`、`streaming?: boolean`

## ShareDialog

- **路径**: `features/tuyile/components/ShareDialog.vue`
- **职责**: 分享弹窗，流程为确认（条数统计 / 预览 / 裁剪说明）→ 成功（链接 / 复制 / 打开），微信内额外提供引导蒙层
- **Props**: `open: boolean`、`bubbles: ChatBubble[]`
- **Emits**: `update:open(boolean)`
