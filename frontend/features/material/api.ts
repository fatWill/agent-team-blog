import type { MaterialItem } from './types'

/** 获取所有材料条目列表 */
export async function apiFetchMaterials(): Promise<MaterialItem[]> {
  const res = await $fetch<{ list: MaterialItem[] }>('/api/materials')
  return res.list ?? []
}

/** 获取单个材料条目详情 */
export async function apiFetchMaterial(id: number): Promise<MaterialItem | null> {
  try {
    return await $fetch<MaterialItem>(`/api/materials/${id}`)
  } catch {
    return null
  }
}

/** 新建材料条目 */
export async function apiCreateMaterial(
  data: Omit<MaterialItem, 'id' | 'createdAt' | 'updatedAt'>,
): Promise<MaterialItem> {
  return $fetch<MaterialItem>('/api/materials', { method: 'POST', body: data })
}

/** 更新材料条目 */
export async function apiUpdateMaterial(
  id: number,
  data: Partial<Omit<MaterialItem, 'id' | 'createdAt' | 'updatedAt'>>,
): Promise<MaterialItem> {
  return $fetch<MaterialItem>(`/api/materials/${id}`, { method: 'PUT', body: data })
}

/** 删除材料条目 */
export async function apiDeleteMaterial(id: number): Promise<void> {
  await $fetch(`/api/materials/${id}`, { method: 'DELETE' })
}
