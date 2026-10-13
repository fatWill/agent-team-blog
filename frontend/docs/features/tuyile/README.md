# 图一乐（tuyile）模块

## 概述

玩具区第三个独立应用「图一乐」（AI 算命），路由 `/toys/tuyile`。本期包含老黄历（单日完整字段 + 月历切换 + 十二时辰）和基于 AG-UI 协议的运势聊天；「星座运势」「生辰八字」入口显示「即将上线」。

后端接口文档：`backend/docs/api/tuyile.md`（后端 commit `9205fc6`）。

## 目录结构

```
features/tuyile/
├── types.ts                       # AlmanacDay / AlmanacMonth / AG-UI 事件与消息类型
├── api.ts                         # 黄历 $fetch + AG-UI SSE 流式请求与解析
├── utils.ts                       # 北京时间日期工具、时辰换算、Markdown 渲染（禁 HTML）
├── composables/
│   ├── useAlmanacDay.ts           # 单日黄历 useAsyncData（SSR 预取）
│   └── useTuyileChat.ts           # 聊天状态：流式拼接、停止、错误、429/503 冷却
├── components/                    # 见 components.md
└── index.ts                       # 公共导出
server/api/toys/tuyile/            # Nitro 透传：almanac/index.get、almanac/month.get、chat.post
server/utils/tuyileProxy.ts        # JSON 透传（保留状态码/Retry-After/Cache-Control）+ SSE 流式透传
```

## AG-UI 对接

- 未引入 `@ag-ui/client`：它依赖 rxjs/zod 等（体积较大），而后端只用 5 种事件。改为 `fetch` + `ReadableStream` 手写 SSE 解析（`readSseEvents`），SSR 无兼容问题（聊天组件包在 `<ClientOnly>` 里）。
- 请求体 `RunAgentInput {threadId, runId, messages}`，每次都带上完整历史；进入聊天时发 `messages: []`，让图一乐主动打招呼。
- 事件处理：`TEXT_MESSAGE_START` 新建气泡，`TEXT_MESSAGE_CONTENT` 追加 `delta`，`TEXT_MESSAGE_END` 结束光标，`RUN_ERROR` 展示 `message`；`: keep-alive` 注释行和未知事件都忽略。
- 停止生成：`AbortController.abort()`，保留已生成的部分；Nitro 代理在客户端断开后中止上游请求，释放后端并发名额。
- HTTP 错误：400、429、503 读后端 `statusMessage`；429/503 按 `Retry-After` 倒计时并禁用发送。请求失败且没有任何回复时，撤回该条用户消息，把文本还给输入框，避免历史里出现连续两条 user。
- 隐私：对话只存在组件内存，不写 localStorage 或 sessionStorage。

## SSE 缓冲

Nitro 层用 `sendStream` 逐块转发，并设置 `X-Accel-Buffering: no` 和 `Cache-Control: no-cache, no-transform`。本地实测各事件按上游节奏到达，没有被攒批。

## 二期（v2.17.0）

- **思考状态**：`CUSTOM name=tuyile.status` 驱动状态条（`value.text`；`querying` 时附 `range`，同年省略年份，如 `2028-02-02~02-16`）。`REASONING_*` 只识别不渲染；`RUN_STARTED` 时先显示「掐指一算中…」兜底；`TEXT_MESSAGE_START` / `RUN_ERROR` / 流结束时隐藏。回传历史只含 user / assistant。
- **错误码**：`RUN_ERROR`（含 `LLM_TIMEOUT`）直接展示后端 `message`。
- **文案**：placeholder 5 条随机起始、每 6 秒轮换（输入非空时不换）；快捷提问 chip 点击直接发送，生成中 / 冷却中禁用。
- **分享**：`ShareDialog` 先用 `planShare` 按后端限制裁剪（只保留 user/assistant 非空消息，单条 ≤4000 字、最近 50 条、请求体 ≤120KB），在弹窗中说明裁剪结果，再 `POST /share`；400/413/429 显示后端文案，429 按 `Retry-After` 倒计时。成功后显示链接，可复制（Clipboard API，失败回退 `execCommand`）或直接打开。
- **微信分支**：UA 含 `MicroMessenger` 时，拿到 url 后调用 `setTuyileWechatCard`（内部 `initWechatShare`，签名 URL = 当前地址去掉 `#`、保留 query；失败用首次进入 URL 重签），然后显示右上角「…」引导蒙层。JS-SDK 无法直接拉起转发面板。
- **渲染安全**：`ChatMessage` 统一渲染，assistant 用 markdown-it `html:false`（原始 HTML 转义），user 用文本插值；分享页 SEO 文本的尖括号转为全角。
