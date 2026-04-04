# 鉴权领域

## 概述

鉴权系统负责管理员登录、Token 管理和路由保护。采用 Cookie Token 方案，支持 72 小时滚动续期，SSR 阶段正确转发 Cookie。

## 当前代码位置

| 功能 | 文件 | 说明 |
|------|------|------|
| 登录页 UI | `pages/login.vue` | 账号密码表单，登录成功跳转 redirect 或 /home |
| 认证状态 | `stores/auth.ts` | Pinia store：`isLoggedIn`、`username`、`checkAuth()`、`logout()` |
| 路由守卫 | `middleware/auth.ts` | 保护 `/admin` 路由，未登录跳转 `/login?redirect=...` |
| API 封装 | `utils/api.ts` | `apiLogin()` |
| 服务端 - 登录 | `server/api/auth/login.post.ts` | 账号密码验证 + 生成 Token + 写入 Cookie |
| 服务端 - 检查 | `server/api/auth/check.get.ts` | 验证 Cookie Token 有效性 + 滚动续期 |
| 服务端 - 登出 | `server/api/auth/logout.post.ts` | 清除 Token + 删除 Cookie |
| 服务端工具 | `server/utils/auth.ts` | Token 生成/验证/续期/鉴权中间件 |

## 鉴权流程

### 登录流程

```
1. 用户输入账号密码
2. → POST /api/auth/login
3. → server/utils/auth.ts 验证密码（bcrypt）
4. → 生成随机 Token，写入 MySQL auth_tokens 表
5. → 设置 httpOnly Cookie（auth_token）
6. → 前端 authStore.setLoggedIn(true, username)
7. → router.push(redirect || '/home')
```

### 路由守卫流程（middleware/auth.ts）

```
1. 访问 /admin 路由
2. → 非 /admin 开头？直接放行
3. → 客户端 + authStore.isLoggedIn？直接放行（内存态）
4. → 客户端 + 无 auth_token Cookie？跳转登录页（省去网络请求）
5. → authStore.checkAuth()
6.   → SSR 阶段：useRequestHeaders(['cookie']) 转发浏览器 Cookie
7.   → GET /api/auth/check 验证 Token
8.   → 成功：更新 lastActiveAt（滚动续期），返回 username
9.   → 失败：跳转 /login?redirect=/admin
```

### Token 续期机制

- Token 存储在 MySQL `auth_tokens` 表
- 字段：`token`、`username`、`lastActiveAt`、`createdAt`
- TTL：72 小时（基于 `lastActiveAt`）
- 每次 `check` 成功时更新 `lastActiveAt`，实现滚动续期
- 服务重启不丢失登录态（Token 持久化在 MySQL）

### SSR Cookie 转发

关键代码在 `stores/auth.ts` 的 `checkAuth()` 方法中：

```typescript
const headers = import.meta.server ? useRequestHeaders(['cookie']) : {}
const res = await $fetch('/api/auth/check', { headers })
```

- **SSR 阶段**：Nuxt 服务端 `$fetch` 不会自动携带浏览器 Cookie，需要通过 `useRequestHeaders` 手动转发
- **客户端阶段**：浏览器 `$fetch` 自动携带 Cookie，传空对象即可

## 安全设计

| 措施 | 说明 |
|------|------|
| httpOnly Cookie | Token 不可被 JavaScript 读取 |
| bcrypt 密码哈希 | 密码不明文存储 |
| Token 随机生成 | UUID v4，不可预测 |
| 72h 自动过期 | 长时间不活跃自动登出 |
| 服务端鉴权 | 所有写操作 API 都通过 `server/utils/auth.ts` 验证 Token |
