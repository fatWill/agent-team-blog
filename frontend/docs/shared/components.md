# 通用组件清单

## AppLoading

- **路径**：`components/AppLoading.vue`
- **职责**：跨端 Loading 加载组件
- **Props**：

| Prop | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| tip | string | ❌ | — | 加载提示文字 |

- **行为**：
  - PC 端（≥768px）：使用 `ant-design-vue` 的 `<ASpin size="large" :tip="tip" />`
  - 移动端（<768px）：使用自定义 SVG spinner + 可选提示文字
  - 通过 `useDevice()` 的 `isMobile` 判断端类型

- **使用示例**：

```vue
<AppLoading tip="加载中..." />
```

- **样式**：
  - 居中显示，`py-20` 上下内边距
  - 移动端 spinner：蓝色（#1677ff），0.8s 旋转动画
  - 支持 Dark 模式

---

> ⚠️ 当前项目仅有 1 个全局组件。其他 UI 功能通过 `utils/ui.ts`（命令式 Toast/Dialog）和 `utils/mobileUI.ts`（移动端命令式组件）实现。
