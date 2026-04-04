# 通用 Composables 清单

## useTheme

- **路径**：`composables/useTheme.ts`
- **职责**：Dark/Light 主题切换，支持 SSR 零闪烁

### 返回值

| 属性/方法 | 类型 | 说明 |
|----------|------|------|
| `colorMode` | `Ref<'light' \| 'dark'>` | 当前主题模式 |
| `isDark` | `ComputedRef<boolean>` | 是否为 Dark 模式 |
| `toggleTheme()` | `() => Promise<void>` | 切换主题（写 Cookie + DOM + Redis） |
| `initTheme()` | `() => void` | 初始化主题到 DOM（仅客户端） |

### 实现细节

- **持久化**：Cookie `color-mode`（30 天 maxAge）+ Redis 异步存储
- **默认值**：`dark`（首次访问默认 Dark 模式）
- **SSR 支持**：通过 `useCookie` 在服务端读取 Cookie 值，`app.vue` 中用 `useHead` 注入 `<html class="dark">`
- **切换流程**：更新 `colorMode` → 写 Cookie → 应用到 DOM → 异步存入 Redis

### 使用示例

```typescript
const { isDark, toggleTheme, initTheme } = useTheme()

// 在 layout 中初始化
onMounted(() => initTheme())

// 在按钮中切换
<button @click="toggleTheme">{{ isDark ? '☀️' : '🌙' }}</button>
```

---

## useDevice

- **路径**：`composables/useDevice.ts`
- **职责**：设备类型检测（移动端/PC 端）

### 返回值（useDevice 组合式函数）

| 属性 | 类型 | 说明 |
|------|------|------|
| `isMobile` | `Ref<boolean>` | 是否为移动端（<768px） |

### 导出函数（命令式调用）

| 函数 | 签名 | 说明 |
|------|------|------|
| `isMobileDevice()` | `() => boolean` | 非响应式移动端检测，用于命令式场景（如 `showSuccess()`） |

### 实现细节

- 断点：768px（与 Tailwind `md` 断点一致）
- 响应式：监听 `resize` 事件实时更新
- SSR 安全：服务端返回 `false`

### 使用示例

```typescript
// 组件内（响应式）
const { isMobile } = useDevice()

// 工具函数内（命令式）
import { isMobileDevice } from '~/composables/useDevice'
if (isMobileDevice()) {
  MobileToast.success('操作成功')
}
```
