/**
 * GET /api/admin/wechat/server-ip
 * 获取 VPS 出口公网 IP（需鉴权）
 * 字段映射：后端 outboundIP/hint/mpConsoleURL → 前端 WechatServerIp
 */
export default defineEventHandler(async (event) => {
  const data = await proxyToBackend(event, '/api/admin/wechat/server-ip')
  return {
    ip: data.outboundIP,
    hint: data.hint,
    mpConsoleURL: data.mpConsoleURL,
  }
})
