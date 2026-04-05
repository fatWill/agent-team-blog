import type { H3Event } from 'h3'

const BACKEND_BASE = process.env.BACKEND_URL || 'http://127.0.0.1:8080'

/**
 * 通用代理：将请求转发到 Go 后端
 * 自动透传 Cookie 和 Set-Cookie
 */
export async function proxyToBackend(
  event: H3Event,
  path: string,
  options?: {
    method?: string
    body?: any
    headers?: Record<string, string>
    query?: Record<string, string>
  },
): Promise<any> {
  const cookie = getRequestHeader(event, 'cookie') || ''

  // 构建 URL（含 query 参数）
  let url = `${BACKEND_BASE}${path}`
  if (options?.query) {
    const qs = new URLSearchParams(options.query).toString()
    if (qs) url += `?${qs}`
  }

  const fetchOptions: RequestInit = {
    method: options?.method || 'GET',
    headers: {
      'Content-Type': 'application/json',
      cookie,
      ...getIpHeaders(event),
      ...(options?.headers || {}),
    },
  }

  if (options?.body !== undefined) {
    fetchOptions.body = typeof options.body === 'string' ? options.body : JSON.stringify(options.body)
  }

  const res = await fetch(url, fetchOptions)

  // 透传 Set-Cookie（登录/续期/theme 等场景）
  const setCookieHeader = res.headers.get('set-cookie')
  if (setCookieHeader) {
    setResponseHeader(event, 'set-cookie', setCookieHeader)
  }

  if (!res.ok) {
    const errBody = await res.json().catch(() => ({ error: res.statusText }))
    throw createError({
      statusCode: res.status,
      statusMessage: errBody?.message || errBody?.error || res.statusText,
    })
  }

  return res.json()
}

/**
 * 透传 multipart/form-data 原始请求到 Go 后端
 * 用于文件上传场景（upload、photos.post）
 */
export async function proxyRawToBackend(
  event: H3Event,
  path: string,
  options?: {
    method?: string
  },
): Promise<any> {
  const cookie = getRequestHeader(event, 'cookie') || ''
  const contentType = getRequestHeader(event, 'content-type') || ''

  // 读取原始请求体
  const rawBody = await readRawBody(event, false)

  const url = `${BACKEND_BASE}${path}`

  const res = await fetch(url, {
    method: options?.method || 'POST',
    headers: {
      'content-type': contentType,
      cookie,
      ...getIpHeaders(event),
    },
    body: rawBody,
  })

  // 透传 Set-Cookie
  const setCookieHeader = res.headers.get('set-cookie')
  if (setCookieHeader) {
    setResponseHeader(event, 'set-cookie', setCookieHeader)
  }

  if (!res.ok) {
    const errBody = await res.json().catch(() => ({ error: res.statusText }))
    throw createError({
      statusCode: res.status,
      statusMessage: errBody?.message || errBody?.error || res.statusText,
    })
  }

  return res.json()
}

/**
 * 获取客户端真实 IP 并构建透传 header
 */
function getIpHeaders(event: H3Event): Record<string, string> {
  const headers: Record<string, string> = {}

  const xff = getRequestHeader(event, 'x-forwarded-for')
  if (xff) headers['x-forwarded-for'] = xff

  const xri = getRequestHeader(event, 'x-real-ip')
  if (xri) headers['x-real-ip'] = xri

  // 如果没有 X-Forwarded-For，用 socket 的 remoteAddress 填充
  if (!xff) {
    const remoteAddr = event.node?.req?.socket?.remoteAddress
    if (remoteAddr) {
      headers['x-forwarded-for'] = remoteAddr
      headers['x-real-ip'] = remoteAddr
    }
  }

  return headers
}
