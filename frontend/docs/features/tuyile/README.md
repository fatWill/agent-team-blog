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
