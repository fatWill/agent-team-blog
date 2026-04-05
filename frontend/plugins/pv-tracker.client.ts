/**
 * PV/UV 自动上报插件（仅客户端执行）
 * 每次路由变化时静默上报访问记录
 */
export default defineNuxtPlugin(() => {
  const router = useRouter()

  // 获取或生成设备唯一 ID
  function getDeviceId(): string {
    let id = localStorage.getItem('_did')
    if (!id) {
      id = crypto.randomUUID()
      localStorage.setItem('_did', id)
    }
    return id
  }

  // 上报函数（fire and forget，不阻塞）
  function track(path: string) {
    const deviceId = getDeviceId()
    const referer = document.referrer || ''
    $fetch('/api/pv/record', {
      method: 'POST',
      body: { path, device_id: deviceId, referer },
    }).catch(() => {}) // 静默失败，不影响用户体验
  }

  // 初始上报（首次进入）
  track(router.currentRoute.value.path)

  // 路由变化时上报
  router.afterEach((to) => {
    track(to.path)
  })
})
