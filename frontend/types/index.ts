/**
 * 全局类型 re-export（向后兼容层）
 *
 * 实际类型定义已按领域拆分到各 feature 模块中：
 * - features/article/types.ts  — 文章相关类型
 * - features/album/types.ts    — 相册/照片相关类型
 * - features/auth/types.ts     — 鉴权相关类型
 * - features/guestbook/types.ts — 留言板相关类型
 * - features/changelog/types.ts — 更新日志相关类型
 * - shared/types.ts            — 通用类型（NavItem、TabItem、Profile）
 *
 * 新代码请直接从对应 feature 模块导入，例如：
 *   import type { ArticleDetail } from '~/features/article'
 */

// === 通用类型 ===
export type { NavItem, TabItem, Profile } from '~/shared/types'

// === 文章领域 ===
export type {
  ArticleListItem,
  ArticleDetail,
  ArticleListResponse,
  CreateArticleRequest,
  UpdateArticleRequest,
} from '~/features/article'

// === 相册领域 ===
export type {
  AlbumItem,
  AlbumListResponse,
  PhotoItem,
  PhotoListResponse,
} from '~/features/album'

// === 鉴权领域 ===
export type {
  LoginRequest,
  LoginResponse,
} from '~/features/auth'

// === 留言板领域 ===
export type {
  MessageItem,
  MessageListResponse,
} from '~/features/guestbook'

// === 更新日志领域 ===
export type {
  ChangelogItem,
  ChangelogResponse,
} from '~/features/changelog'
