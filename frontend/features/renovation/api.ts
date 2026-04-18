import type { RenovationArticle, RenovationArticleListItem } from './types'

/** 获取装修知识文章列表 */
export async function apiFetchRenovationArticles(): Promise<RenovationArticleListItem[]> {
  const res = await $fetch<{ list: RenovationArticleListItem[] }>('/api/renovation/articles')
  return res.list ?? []
}

/** 获取装修知识文章详情 */
export async function apiFetchRenovationArticle(id: string | number): Promise<RenovationArticle | null> {
  try {
    return await $fetch<RenovationArticle>(`/api/renovation/articles/${id}`)
  } catch {
    return null
  }
}
