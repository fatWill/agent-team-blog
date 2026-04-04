import type { LoginRequest, LoginResponse } from './types'

/** 登录 */
export async function apiLogin(data: LoginRequest): Promise<LoginResponse> {
  const res = await $fetch<LoginResponse>('/api/auth/login', {
    method: 'POST',
    body: data,
  })
  return res
}
