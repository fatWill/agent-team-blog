import type { H3Event } from 'h3'

/**
 * 玩具「图一乐」专用代理工具
 * - proxyTuyileJson：GET 透传，保留后端状态码、错误体 {error,statusCode,statusMessage} 与 Cache-Control
 * - proxyTuyileStream：POST AG-UI SSE 透传，逐块转发不缓冲，客户端断开时中止上游
 */

function backendBase(): string {
  return process.env.BACKEND_URL || 'http://127.0.0.1:8080'
}

function forwardIpHeaders(event: H3Event): Record<string, string> {
  const headers: Record<string, string> = {}
  const xff = getRequestHeader(event, 'x-forwarded-for')
  const xri = getRequestHeader(event, 'x-real-ip')
  if (xff) headers['x-forwarded-for'] = xff
  if (xri) headers['x-real-ip'] = xri
  if (!xff) {
    const remote = event.node?.req?.socket?.remoteAddress
    if (remote) {
      headers['x-forwarded-for'] = remote
      if (!xri) headers['x-real-ip'] = remote
    }
  }
  return headers
}

const PASS_HEADERS = ['cache-control', 'retry-after']

function passHeaders(event: H3Event, res: Response) {
  for (const name of PASS_HEADERS) {
    const v = res.headers.get(name)
    if (v) setResponseHeader(event, name, v)
  }
}

async function readJsonOrError(res: Response): Promise<unknown> {
  const text = await res.text()
  try {
    return JSON.parse(text)
  } catch {
    return { error: true, statusCode: res.status, statusMessage: res.statusText || '服务暂不可用' }
  }
}

export async function proxyTuyileJson(event: H3Event, path: string, method: 'GET' | 'POST' = 'GET'): Promise<unknown> {
  const query = getQuery(event)
  const qs = new URLSearchParams()
  for (const [k, v] of Object.entries(query)) {
    if (typeof v === 'string') qs.set(k, v)
  }
  const url = `${backendBase()}${path}${qs.toString() ? `?${qs.toString()}` : ''}`
  const init: RequestInit = { method, headers: { accept: 'application/json', ...forwardIpHeaders(event) } }
  if (method === 'POST') {
    // 原样透传请求体（后端自行做 JSON / 128KB 校验，保留 400/413 语义）
    const body = await readRawBody(event, false)
    init.headers = { ...init.headers, 'content-type': 'application/json' }
    init.body = body ? new Uint8Array(body) : undefined
  }
  let res: Response
  try {
    res = await fetch(url, init)
  } catch {
    setResponseStatus(event, 502)
    return { error: true, statusCode: 502, statusMessage: '后端服务不可用' }
  }
  passHeaders(event, res)
  setResponseStatus(event, res.status)
  return readJsonOrError(res)
}

export async function proxyTuyileStream(event: H3Event, path: string) {
  const body = await readRawBody(event, false)
  const controller = new AbortController()
  // 客户端断开（停止生成 / 关闭页面）时中止上游，释放后端并发名额
  event.node.req.on('close', () => controller.abort())
  event.node.res.on('close', () => controller.abort())

  let res: Response
  try {
    res = await fetch(`${backendBase()}${path}`, {
      method: 'POST',
      headers: {
        'content-type': 'application/json',
        accept: 'text/event-stream',
        ...forwardIpHeaders(event),
      },
      body: body ? new Uint8Array(body) : undefined,
      signal: controller.signal,
    })
  } catch {
    setResponseStatus(event, 502)
    return { error: true, statusCode: 502, statusMessage: '后端服务不可用' }
  }

  passHeaders(event, res)
  setResponseStatus(event, res.status)

  const contentType = res.headers.get('content-type') || ''
  if (!res.ok || !contentType.includes('text/event-stream') || !res.body) {
    return readJsonOrError(res)
  }

  setResponseHeaders(event, {
    'content-type': 'text/event-stream; charset=utf-8',
    'cache-control': 'no-cache, no-transform',
    'x-accel-buffering': 'no',
    connection: 'keep-alive',
  })
  return sendStream(event, res.body)
}
