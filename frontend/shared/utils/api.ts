import type { Profile } from '~/shared/types'

/** 获取博主个人资料 */
export async function apiGetProfile(): Promise<Profile> {
  const res = await $fetch<Profile>('/api/profile')
  return res
}

/** 更新博主个人资料 */
export async function apiUpdateProfile(data: Partial<Profile>): Promise<Profile> {
  const res = await $fetch<Profile>('/api/profile', {
    method: 'PUT',
    body: data,
  })
  return res
}

/** 上传图片 */
export async function apiUploadImage(file: File): Promise<{ url: string }> {
  const formData = new FormData()
  formData.append('file', file)
  const res = await $fetch<{ url: string }>('/api/upload', {
    method: 'POST',
    body: formData,
  })
  return res
}
