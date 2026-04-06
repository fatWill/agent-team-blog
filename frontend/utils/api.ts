/**
 * API 函数 re-export（向后兼容层）
 *
 * 实际 API 函数已按领域拆分到各 feature 模块中：
 * - features/article/api.ts   — 文章相关 API
 * - features/album/api.ts     — 相册/照片相关 API
 * - features/auth/api.ts      — 鉴权相关 API
 * - features/guestbook/api.ts — 留言板相关 API
 * - shared/utils/api.ts       — 通用 API（profile、upload）
 *
 * 新代码请直接从对应 feature 模块导入，例如：
 *   import { apiFetchArticles } from '~/features/article'
 */

// === 鉴权领域 ===
export { apiLogin } from '~/features/auth'

// === 文章领域 ===
export {
  apiFetchArticles,
  apiFetchArticle,
  apiCreateArticle,
  apiUpdateArticle,
  apiDeleteArticle,
  apiGetArticles,
  apiToggleArticleLike,
  apiGetArticleLikeStatus,
  apiGetArticleLikeStatusBatch,
  apiRecordArticleView,
  apiGetRandomArticle,
  apiSearchArticles,
} from '~/features/article'

// === 相册领域 ===
export {
  apiGetAlbums,
  apiCreateAlbum,
  apiUpdateAlbum,
  apiDeleteAlbum,
  apiGetPhotos,
  apiAddPhoto,
  apiDeletePhoto,
  apiVerifyAlbumPassword,
  apiVerifyPhotoPassword,
  apiUpdatePhoto,
} from '~/features/album'

// === 留言板领域 ===
export {
  apiGetMessages,
  apiCreateMessage,
  apiDeleteMessage,
} from '~/features/guestbook'

// === 通用（profile / upload） ===
export {
  apiGetProfile,
  apiUpdateProfile,
  apiUploadImage,
} from '~/shared/utils/api'
