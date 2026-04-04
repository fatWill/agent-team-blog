# 留言板领域

## 概述

留言板是博客主页的一个 Tab，允许访客留言互动。支持留言发布、编辑（每日一次）、管理员删除等功能。

## 核心实体

### MessageItem

```typescript
interface MessageItem {
  id: number
  nickname: string       // 昵称（可选，默认"匿名"）
  content: string        // 留言内容（必填）
  isOwn: boolean         // 是否是当前设备的留言
  canEdit: boolean       // 今天是否还能编辑
  createdAt: string
  updatedAt: string
}
```

## 当前代码位置

| 功能 | 文件 | 说明 |
|------|------|------|
| 留言板 UI | `pages/home.vue`（留言板 Tab） | 留言输入区 + 留言列表 |
| 留言管理 | `pages/admin.vue`（留言管理 Tab） | 管理员删除留言 |
| API 封装 | `utils/api.ts` | `apiGetMessages()`、`apiCreateMessage()`、`apiUpdateMessage()`、`apiDeleteMessage()` |
| 类型定义 | `types/index.ts` | `MessageItem`、`MessageListResponse` |
| 服务端 API | `server/api/messages/` | 3 个接口文件 |
| 服务端 DAO | `server/utils/messages.ts` | MySQL CRUD + deviceId 关联 |

## API 接口

| 方法 | 路径 | 鉴权 | 说明 |
|------|------|------|------|
| GET | `/api/messages` | ❌ | 留言列表（`?deviceId=xxx` 标记自己的留言） |
| POST | `/api/messages` | ❌ | 新增留言（IP 限频） |
| PUT | `/api/messages/:id` | ❌ | 修改留言（每日一次限制） |
| DELETE | `/api/messages/:id` | ✅ | 管理员删除留言 |

## 业务流程

### 发布留言

```
1. 用户输入昵称（可选）+ 内容（必填）
2. → apiCreateMessage({ deviceId, nickname, content })
3. → POST /api/messages
4. → IP 限频检查（rateLimit.ts）
5. → 写入 MySQL messages 表
6. → 刷新留言列表
```

### 编辑留言

```
1. 已有留言的用户，按钮显示为"修改留言"
2. → 点击后进入编辑模式
3. → apiUpdateMessage(id, { deviceId, nickname, content })
4. → PUT /api/messages/:id
5. → 检查每日修改限制（403 今天已修改过）
6. → 更新成功，刷新列表
```

### 设备标识

- 使用 `deviceId` 机制（localStorage `blog_device_id`，UUID v4）
- 标识当前设备的留言，实现 `isOwn` 和 `canEdit` 判断
- 无需登录即可留言

## 限制规则

| 规则 | 说明 |
|------|------|
| IP 限频 | 同一 IP 短时间内不能频繁发布（429 Too Many Requests） |
| 每日修改 | 每条留言每天只能修改一次（403 Forbidden） |
| 内容必填 | content 不能为空 |
| 昵称可选 | nickname 为空时显示为"匿名" |

## 错误处理

| 状态码 | 说明 | 前端处理 |
|--------|------|---------|
| 429 | 操作太频繁 | `showError('操作太频繁，请稍后再试')` |
| 403 | 今天已修改过 | `showError('今天已经修改过了')` |
| 400 | 参数错误 | `showError('请输入留言内容')` |
