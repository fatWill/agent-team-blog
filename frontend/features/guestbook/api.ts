import type { MessageItem, MessageListResponse } from './types'

/** 获取留言列表 */
export async function apiGetMessages(deviceId?: string): Promise<MessageListResponse> {
  const params: Record<string, string> = {}
  if (deviceId) params.deviceId = deviceId
  const res = await $fetch<MessageListResponse>('/api/messages', { params })
  return res
}

/** 新增留言 */
export async function apiCreateMessage(data: { deviceId: string; nickname?: string; content: string }): Promise<MessageItem> {
  const res = await $fetch<MessageItem>('/api/messages', {
    method: 'POST',
    body: data,
  })
  return res
}

/** 修改留言 */
export async function apiUpdateMessage(id: number, data: { deviceId: string; nickname?: string; content: string }): Promise<MessageItem> {
  const res = await $fetch<MessageItem>(`/api/messages/${id}`, {
    method: 'PUT',
    body: data,
  })
  return res
}

/** 删除留言（管理员） */
export async function apiDeleteMessage(id: number): Promise<{ success: boolean }> {
  const res = await $fetch<{ success: boolean }>(`/api/messages/${id}`, {
    method: 'DELETE',
  })
  return res
}
