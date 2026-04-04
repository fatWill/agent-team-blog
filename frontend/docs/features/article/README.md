# 文章领域

## 概述

文章是博客的核心内容。支持文章的创建、编辑、删除、列表展示、详情阅读和点赞交互。

## 核心实体

### ArticleListItem

文章列表项（不含正文内容），用于列表展示。

```typescript
interface ArticleListItem {
  id: string              // UUID
  title: string           // 标题
  summary: string         // 摘要
  coverImage: string      // 封面图 URL
  likeCount: number       // 点赞数
  createdAt: string       // 创建时间
  updatedAt: string       // 更新时间
}
```

### ArticleDetail

文章详情（继承 ArticleListItem），包含 Tiptap JSON 格式的正文内容。

```typescript
interface ArticleDetail extends ArticleListItem {
  content: Record<string, unknown>  // Tiptap JSON
}
```

## 当前代码位置

> ⚠️ 文章相关代码分散在多个文件中，尚未拆分为独立 feature 模块。

| 功能 | 文件 | 说明 |
|------|------|------|
| 文章列表 UI | `pages/home.vue`（文章 Tab） | 虚拟滚动列表，卡片展示封面/标题/摘要/点赞数 |
| 文章详情页 | `pages/articles/[id].vue` | Tiptap 只读渲染 + 底部点赞按钮 |
| 文章管理 | `pages/admin.vue`（写文章 Tab / 文章管理 Tab） | Tiptap 编辑器 + 文章列表管理 |
| API 封装 | `utils/api.ts` | `apiFetchArticles()`、`apiFetchArticle()`、`apiCreateArticle()`、`apiUpdateArticle()`、`apiDeleteArticle()`、`apiToggleArticleLike()`、`apiGetArticleLikeStatus()`、`apiGetArticleLikeStatusBatch()` |
| 类型定义 | `types/index.ts` | `ArticleListItem`、`ArticleDetail`、`CreateArticleRequest`、`UpdateArticleRequest`、`ArticleListResponse` |
| 服务端 API | `server/api/articles/` | 6 个接口文件 |
| 服务端 DAO | `server/utils/articles.ts` | MySQL CRUD + like_count 冗余字段 |

## API 接口

| 方法 | 路径 | 鉴权 | 说明 |
|------|------|------|------|
| GET | `/api/articles` | ❌ | 文章列表（支持 `?title=xxx` 关键词筛选） |
| GET | `/api/articles/:id` | ❌ | 文章详情 |
| POST | `/api/articles` | ✅ | 创建文章 |
| PUT | `/api/articles/:id` | ✅ | 更新文章（部分更新） |
| DELETE | `/api/articles/:id` | ✅ | 物理删除 |
| POST | `/api/articles/:id/like` | ❌ | 点赞/取消（deviceId 切换） |
| GET | `/api/articles/:id/like-status` | ❌ | 查询单篇点赞状态 |
| GET | `/api/articles/like-status-batch` | ❌ | 批量查询点赞状态（ids 逗号分隔） |

## 业务流程

### 文章列表加载（SSR）

```
1. home.vue onMounted / useAsyncData
2. → apiFetchArticles()
3. → GET /api/articles
4. → server/utils/articles.ts 查询 MySQL
5. → SSR 预取注入 HTML
6. → 客户端 hydration，DynamicScroller 虚拟滚动渲染
```

### 文章点赞

```
1. 用户点击爱心按钮
2. → 乐观更新 UI（立即切换动画）
3. → apiToggleArticleLike(articleId, deviceId)
4. → AbortController 取消旧请求（快速切换场景）
5. → POST /api/articles/:id/like
6. → article_likes 表 INSERT/DELETE + articles 表 like_count 更新
7. → 失败时回滚 UI（区分 AbortError 和真正网络错误）
```

### 批量点赞状态查询

首页加载文章列表后，一次性批量查询所有文章的点赞状态，避免 N+1 请求：

```
GET /api/articles/like-status-batch?ids=id1,id2,id3&deviceId=xxx
→ 返回 { likedIds: ['id1', 'id3'] }
```

## 富文本编辑器

使用 Tiptap 作为富文本编辑器：

- **编辑模式**：`admin.vue` 中的写文章 Tab，支持标题、正文、图片插入、链接
- **只读模式**：`articles/[id].vue` 中 Tiptap 渲染 JSON content
- **依赖**：`@tiptap/vue-3`、`@tiptap/starter-kit`、`@tiptap/extension-image`、`@tiptap/extension-link`、`@tiptap/pm`
