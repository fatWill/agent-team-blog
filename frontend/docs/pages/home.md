# 酷炫首页（/）

## 基本信息

| 属性 | 值 |
|------|------|
| 路由 | `/` |
| 文件 | `pages/index.vue`（17KB） |
| 渲染模式 | CSR（纯 Canvas 交互，无 SEO 需求） |
| 依赖 | `@chenglou/pretext` |

## 功能描述

酷炫欢迎页面，作为博客的入口。使用 Canvas 绘制代码文字滚动背景，中央放置可拖拽的「进入博客」按钮。

## 核心实现

### Canvas 代码滚动背景

- 使用 `@chenglou/pretext` 库进行文字排版和宽度测量
- 18 行代码片段以不同速度/方向横向无限滚动
- 奇偶行方向相反，产生视差感
- 速度范围：0.4~1.2px/帧
- 每行由 2-3 个代码片段用 `//` 拼接后重复 4 次

### 绕排逻辑

- 检测每行文字与中央按钮的垂直重叠
- 重叠行分为左右两个 clip 区域绘制连续滚动文字
- 按钮拖拽时文字动态绕排

### 响应式参数

| 参数 | 移动端（<768px） | 桌面端（≥768px） |
|------|-----------------|-----------------|
| 字号 | 10-11px | 13-16px |
| 行间距 | 36px | 26px |

### Canvas 渲染优化

- `devicePixelRatio` 高清渲染
- `ctx.setTransform(dpr, 0, 0, dpr, 0, 0)` 每帧重置变换矩阵（修复 scale 累积 bug）
- 按钮位置使用缓存变量（`cachedBtnL/T/R/B`），仅拖拽/初始化/resize 时更新（避免 rAF 内 `getBoundingClientRect` 回流）
- `canvas.measureText()` 测量真实像素宽度（替代 `layoutNextLine` 只测第一逻辑行的问题）

### 拖拽交互

- 移动距离阈值 8px 触发拖拽模式
- 短按（未达阈值）跳转 `/home`
- 拖拽过程中按钮位置实时更新，文字动态绕排

### 按钮入场动画

- `scale(0.3)` 弹性放大到 `scale(1)`
- 过渡曲线：`cubic-bezier(0.34, 1.56, 0.64, 1)`
- 初始定位：`opacity: 0` + `pointer-events: none` 消除占位，`estimateBtnSize` 估算尺寸居中

## 已知 Bug 修复历史

1. **ctx.scale 累积**：改用 `setTransform` 每帧重置
2. **rAF 内 getBoundingClientRect 回流**：改用缓存变量
3. **文字重叠**：`layoutNextLine` 只测第一逻辑行宽度，改用 `canvas.measureText()` 测量
4. **移动端行重叠**：响应式字号/行间距 + resize 时完全重建行数据
5. **拖拽延迟**：移除 setTimeout，改为移动距离阈值
