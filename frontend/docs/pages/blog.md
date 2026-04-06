# 博客主页（/home）

## 基本信息

| 属性 | 值 |
|------|------|
| 路由 | `/home` |
| 文件 | `pages/home.vue`（82KB，大文件） |
| 渲染模式 | SSR + CSR（首屏数据服务端预取） |
| 依赖 | useTheme, useDevice, api.ts, types, vue-virtual-scroller, ant-design-vue |

## 功能描述

博客的主要内容页面，包含个人信息区域和 6 个 Tab 导航。

## 页面结构

### 个人信息区域

- 头像（CDN URL）
- 名字：fatwill
- 简介（从 `/api/profile` 获取）
- GitHub 链接图标

### Tab 导航

| Tab | Emoji | 功能 |
|-----|-------|------|
| 文章 | 📝 | 文章列表（虚拟滚动） |
| 生活 | 🌈 | 相册集 → 照片浏览 → 灯箱预览 |
| 玩具 | 🎮 | 预留 |
| Agent Team | 🤖 | 4 个 AI Agent 卡片展示 |
| 更新日志 | 📋 | Git-log 风格时间轴 |
| 留言板 | 💬 | 留言输入 + 列表 |

### 响应式布局

| 元素 | 移动端（<768px） | PC 端（≥768px） |
|------|-----------------|-----------------|
| 导航 | 左上角汉堡菜单 ☰ + 侧滑抽屉 | 右上角固定按钮组（GitHub + 主题切换 + 管理后台） |
| Tab | 横向可滚动 | 横向排列 |
| 文章列表 | 单列 | 单列（max-width 限制） |
| 相册网格 | 2 列 | 5 列 |
| Agent 卡片 | 1 列 | 2 列 |

### 抽屉导航（移动端）

- 宽度：288px
- 动画：`cubic-bezier(0.16, 1, 0.3, 1)`
- 内容：头像/名称/简介/Tab 导航/GitHub 链接
- 实现：Vue Transition + Teleport to body

## SSR 数据预取

以下数据在服务端预取（`useAsyncData`）：

| 数据 | 函数 | API |
|------|------|-----|
| 文章列表 | `fetchArticles()` | `GET /api/articles` |
| 个人资料 | `fetchProfile()` | `GET /api/profile` |
| 更新日志 | `fetchChangelog()` | `GET /api/changelog` |

以下数据保持懒加载（客户端获取）：

| 数据 | 触发时机 |
|------|---------|
| 相册列表 | 切换到「生活」Tab |
| 照片列表 | 点击相册集 |
| 留言列表 | 切换到「留言板」Tab |
| 文章点赞状态 | 客户端 hydration 后批量查询 |

## 文章 Tab 详情

- **虚拟滚动**：`DynamicScroller` + `DynamicScrollerItem`（`<ClientOnly>` 包裹）
- **SSR Fallback**：静态 `v-for` 渲染前几篇文章
- **卡片内容**：封面图（CDN URL）、标题、摘要、点赞数、发布时间
- **点赞**：爱心动画 + 乐观更新 + AbortController 防快速切换

## 生活 Tab 详情

- **相册集网格**：封面 + 名称 + 照片数 + 锁图标
- **照片浏览**：按年月分组 + 虚拟滚动
- **灯箱**：全屏预览 + 滑动/键盘切换 + 双指缩放 + 点赞/踩
- **密码保护**：模糊遮罩 + 验证弹窗 + sessionStorage 缓存

## Agent Team Tab 详情

- 4 个 AI Agent 卡片：PM、前端、后端、运维
- 每个卡片：图标 + 名称 + 英文名、角色描述、Tags 标签组、模型名徽章
- 不同主题色左边框区分（purple/blue/green/orange）
- 数据写死在前端，无需接口

## 更新日志 Tab 详情

- Git-log 风格时间轴布局
- 版本倒序排列
- 竖线 + 圆点装饰
- 数据来源：`GET /api/changelog`

## 留言板 Tab 详情

- 留言输入区：昵称（可选）+ 内容（必填）+ 发布按钮
- 已有留言时按钮变为「修改留言」
- 留言列表：昵称、内容、相对时间
- 错误处理：429 操作太频繁、403 今天已修改过

## SEO / Meta

```typescript
useHead({
  title: 'fatwill - 个人博客',
  meta: [
    { name: 'description', content: '全栈开发者 fatwill 的个人博客，分享技术与生活。' },
  ],
})
```

## 底部 Footer

```
© 2025 fatwill. All rights reserved. 粤ICP备2025475180号
```

ICP 备案号可点击跳转工信部备案查询页面。
