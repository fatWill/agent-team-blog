# 代码规范与命名约定

## ESLint 配置现状

> ⚠️ 当前项目**未配置 ESLint**，依赖 TypeScript 编译器（`vue-tsc`）进行类型检查和基本代码质量保障。后续可考虑引入 `@nuxt/eslint-config` + `@typescript-eslint`。

### TypeScript 配置

- `tsconfig.json` 继承 Nuxt 自动生成的配置
- `nuxt.config.ts` 中 `typescript.strict: true`
- 构建时通过 `vue-tsc` 进行类型检查

## 代码风格约定

### 缩进与格式

| 规则 | 约定 |
|------|------|
| 缩进 | 2 空格 |
| 引号 | 单引号 `'` |
| 分号 | 无分号 |
| 行尾逗号 | 有（trailing comma） |
| 最大行宽 | 无硬性限制，建议 120 字符 |

### Vue 组件规范

| 规则 | 约定 |
|------|------|
| 组件语法 | `<script setup lang="ts">` 优先 |
| 模板顺序 | `<template>` → `<script setup>` → `<style scoped>` |
| Props 定义 | `defineProps<T>()` 泛型语法 |
| Emits 定义 | `defineEmits<T>()` 泛型语法 |
| 样式作用域 | `<style scoped>` 或 Tailwind 原子类 |

### TypeScript 规范

| 规则 | 约定 |
|------|------|
| `any` 类型 | 禁止使用，必要时用 `unknown` + 类型守卫 |
| 类型导入 | 使用 `import type { ... }` |
| 非空断言 | 尽量避免 `!`，优先使用可选链 `?.` |
| 类型定义位置 | 全局类型放 `types/index.ts`，领域类型放对应模块内 |

## 命名约定

### 文件命名

| 类型 | 规则 | 示例 |
|------|------|------|
| Vue 页面 | kebab-case | `index.vue`、`home.vue`、`login.vue` |
| Vue 组件 | PascalCase | `AppLoading.vue` |
| Composable | camelCase + use 前缀 | `useTheme.ts`、`useDevice.ts` |
| Store | camelCase | `auth.ts`（store 名 `useAuthStore`） |
| 工具函数 | camelCase | `api.ts`、`imageUrl.ts`、`ui.ts` |
| 类型文件 | camelCase | `index.ts` |
| 服务端 API | Nuxt 文件路由约定 | `index.get.ts`、`[id].put.ts` |
| 插件 | kebab-case + 环境后缀 | `antd.client.ts`、`pinia-hydration-fix.ts` |

### 变量/函数命名

| 类型 | 规则 | 示例 |
|------|------|------|
| 组合式函数 | `use` + PascalCase | `useTheme()`、`useDevice()` |
| Store | `use` + PascalCase + `Store` | `useAuthStore` |
| API 函数 | `api` + PascalCase | `apiLogin()`、`apiFetchArticles()` |
| 工具函数 | camelCase | `toCdnUrl()`、`showSuccess()` |
| 常量 | UPPER_SNAKE_CASE | `CHUNK_SIZE`、`CDN_HOST` |
| 接口/类型 | PascalCase | `ArticleDetail`、`LoginRequest` |
| ref 变量 | camelCase | `isLoggedIn`、`colorMode` |
| 事件处理器 | `handle` + 动词 | `handleLogout()`、`handleDelete()` |

### CSS 类名

- **优先使用 Tailwind 原子类**
- 自定义类使用 kebab-case：`mobile-toast-mask`、`mobile-loading-wrap`
- Dark 模式使用 Tailwind `dark:` 前缀

## Tailwind CSS 配置

### 自定义主题色

| 色系 | 用途 | 色值范围 |
|------|------|---------|
| `primary` | 主色调（链接、按钮、激活态） | sky 蓝色系（50-950） |
| `accent` | 强调色（成功、标签） | emerald 绿色系（50-950） |

### 字体栈

- **sans**：Inter → system-ui → PingFang SC → Noto Sans SC
- **mono**：JetBrains Mono → Fira Code → monospace

### 响应式断点

| 断点 | 宽度 | 用途 |
|------|------|------|
| 默认 | < 768px | 移动端（单列布局） |
| `md` | ≥ 768px | 平板/桌面端（多列布局） |
| `lg` | ≥ 1024px | 桌面端（更宽布局） |

## Git 提交规范

### Conventional Commits

格式：`<type>(frontend agent): <描述>`

| type | 用途 |
|------|------|
| `feat` | 新功能 |
| `fix` | 修复 bug |
| `style` | 样式调整（不影响逻辑） |
| `refactor` | 重构 |
| `docs` | 文档变更 |
| `chore` | 构建/工具/依赖变更 |
| `perf` | 性能优化 |

### 提交流程

```bash
git add .
git commit -m "<type>(frontend agent): <描述>"
git push origin main
```
