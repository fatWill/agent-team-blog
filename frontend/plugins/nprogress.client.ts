import NProgress from 'nprogress'

export default defineNuxtPlugin(() => {
  // NProgress 配置
  NProgress.configure({
    showSpinner: false,
    speed: 400,
    minimum: 0.15,
    easing: 'ease',
    trickleSpeed: 200,
  })

  const router = useRouter()

  // 路由导航开始 → 立即启动进度条（比 page:start 更早触发）
  router.beforeEach((_to, _from) => {
    NProgress.start()
  })

  // 路由导航完成 → 完成进度条
  router.afterEach(() => {
    NProgress.done()
  })

  // 路由导航失败 → 也完成进度条
  router.onError(() => {
    NProgress.done()
  })
})
